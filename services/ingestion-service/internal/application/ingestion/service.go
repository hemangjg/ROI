package ingestion

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/ai-finops/ai-finops/packages/shared"
	"github.com/ai-finops/ai-finops/services/ingestion-service/internal/infrastructure/redis"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service implements ingestion flows and outbox writes.
type Service struct {
	pool      *pgxpool.Pool
	queries   *sqlcgen.Queries
	limiter   *redis.RateLimiter
	validate  *validator.Validate
	log       *slog.Logger
}

// NewService constructs an ingestion application service.
func NewService(
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
	limiter *redis.RateLimiter,
	log *slog.Logger,
) *Service {
	return &Service{
		pool:     pool,
		queries:  queries,
		limiter:  limiter,
		validate: validator.New(),
		log:      log,
	}
}

// IngestEvent validates and queues a single usage event via the transactional outbox.
func (s *Service) IngestEvent(ctx context.Context, identity APIKeyIdentity, event UsageEvent) (IngestResult, error) {
	if err := s.limiter.Allow(ctx, identity.KeyID); err != nil {
		return IngestResult{}, err
	}
	if err := s.validateEvent(event); err != nil {
		return IngestResult{}, err
	}

	orgUUID, err := parseUUID(identity.OrgID)
	if err != nil {
		return IngestResult{}, err
	}

	return s.queueEvent(ctx, orgUUID, event)
}

// IngestBatch validates and queues up to 500 usage events.
func (s *Service) IngestBatch(ctx context.Context, identity APIKeyIdentity, events []UsageEvent) (BatchIngestResult, error) {
	if len(events) == 0 {
		return BatchIngestResult{}, shared.BadRequest("invalid_request", "events array is required")
	}
	if len(events) > maxBatchSize {
		return BatchIngestResult{}, shared.BadRequest("batch_too_large", "batch cannot exceed 500 events")
	}

	result := BatchIngestResult{}
	for i, event := range events {
		if _, err := s.IngestEvent(ctx, identity, event); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, batchError(i, err))
			continue
		}
		result.Accepted++
	}
	return result, nil
}

// GetEventStatus returns the current processing status for an ingested event.
func (s *Service) GetEventStatus(ctx context.Context, identity APIKeyIdentity, eventID string) (EventStatus, error) {
	orgUUID, err := parseUUID(identity.OrgID)
	if err != nil {
		return EventStatus{}, err
	}
	eventUUID, err := parseUUID(eventID)
	if err != nil {
		return EventStatus{}, err
	}

	row, err := s.queries.GetOutboxEventByID(ctx, sqlcgen.GetOutboxEventByIDParams{
		ID:    eventUUID,
		OrgID: orgUUID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return EventStatus{}, shared.NotFound("event_not_found", "event not found")
		}
		return EventStatus{}, shared.Internal("lookup outbox event")
	}

	status := statusQueued
	if row.PublishedAt.Valid {
		status = statusQueued
	}

	return EventStatus{
		EventID: pgUUIDString(row.ID),
		Status:  status,
	}, nil
}

// IngestEventForOrg queues an event for an explicit org ID (internal gRPC use).
func (s *Service) IngestEventForOrg(ctx context.Context, orgID string, event UsageEvent) (IngestResult, error) {
	if err := s.validateEvent(event); err != nil {
		return IngestResult{}, err
	}
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return IngestResult{}, err
	}
	return s.queueEvent(ctx, orgUUID, event)
}

func (s *Service) queueEvent(ctx context.Context, orgID pgtype.UUID, event UsageEvent) (IngestResult, error) {
	eventID := uuid.New()
	receivedAt := time.Now().UTC()
	event = normalizeEvent(event)

	payload := OutboxPayload{
		EventID:    eventID.String(),
		OrgID:      pgUUIDString(orgID),
		Event:      event,
		ReceivedAt: receivedAt,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return IngestResult{}, shared.Internal("marshal outbox payload")
	}

	_, err = s.queries.InsertOutboxEvent(ctx, sqlcgen.InsertOutboxEventParams{
		ID:        pgtype.UUID{Bytes: eventID, Valid: true},
		OrgID:     orgID,
		EventType: outboxEventType,
		Payload:   payloadBytes,
	})
	if err != nil {
		return IngestResult{}, shared.Internal("insert outbox event")
	}

	return IngestResult{
		EventID: eventID.String(),
		Status:  statusQueued,
	}, nil
}

func (s *Service) validateEvent(event UsageEvent) error {
	event.IdempotencyKey = strings.TrimSpace(event.IdempotencyKey)
	event.Provider = strings.TrimSpace(strings.ToLower(event.Provider))
	event.Model = strings.TrimSpace(event.Model)

	if event.OccurredAt.IsZero() {
		return shared.BadRequest("invalid_request", "occurred_at is required")
	}

	if err := s.validate.Struct(event); err != nil {
		return shared.BadRequest("validation_error", validationMessage(err))
	}
	if len(event.IdempotencyKey) > maxIdempotencyLen {
		return shared.BadRequest("validation_error", "idempotency_key exceeds max length")
	}
	return nil
}

func normalizeEvent(event UsageEvent) UsageEvent {
	if event.OccurredAt.Location() != time.UTC {
		event.OccurredAt = event.OccurredAt.UTC()
	}
	return event
}

func validationMessage(err error) string {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) && len(validationErrs) > 0 {
		first := validationErrs[0]
		return first.Field() + " failed validation on '" + first.Tag() + "'"
	}
	return "invalid request payload"
}

func batchError(index int, err error) BatchIngestError {
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		return BatchIngestError{
			Index:   index,
			Code:    appErr.Code,
			Message: appErr.Message,
		}
	}
	return BatchIngestError{
		Index:   index,
		Code:    "internal_error",
		Message: "internal server error",
	}
}