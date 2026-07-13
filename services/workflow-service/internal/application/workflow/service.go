package workflow

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/ai-finops/ai-finops/packages/shared"
	chinfra "github.com/ai-finops/ai-finops/services/workflow-service/internal/infrastructure/clickhouse"
	pricinginfra "github.com/ai-finops/ai-finops/services/workflow-service/internal/infrastructure/pricing"
	ingestionv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/ingestion/v1"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	statusProcessed = "processed"
	statusDuplicate = "duplicate"
)

// ProcessResult is returned after processing a usage event.
type ProcessResult struct {
	Status  string
	CostUSD string
}

// Service processes usage events from Kafka or gRPC.
type Service struct {
	pool             *pgxpool.Pool
	queries          *sqlcgen.Queries
	pricing          *pricinginfra.Client
	clickhouse       *chinfra.Client
	idempotencyTTL   time.Duration
	log              *slog.Logger
}

// NewService constructs a workflow application service.
func NewService(
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
	pricing *pricinginfra.Client,
	clickhouse *chinfra.Client,
	idempotencyTTL time.Duration,
	log *slog.Logger,
) *Service {
	return &Service{
		pool:           pool,
		queries:        queries,
		pricing:        pricing,
		clickhouse:     clickhouse,
		idempotencyTTL: idempotencyTTL,
		log:            log,
	}
}

// ProcessUsageEventCreated handles a Kafka UsageEventCreated message.
func (s *Service) ProcessUsageEventCreated(ctx context.Context, created *ingestionv1.UsageEventCreated) (ProcessResult, error) {
	if created == nil || created.GetEvent() == nil {
		return ProcessResult{}, shared.BadRequest("invalid_request", "usage event payload is required")
	}
	return s.process(ctx, created.GetEventId(), created.GetOrgId(), created.GetEvent(), created.GetReceivedAt().AsTime())
}

// ProcessUsageEvent handles workflow.v1.ProcessUsageEventRequest data.
func (s *Service) ProcessUsageEvent(
	ctx context.Context,
	eventID, orgID string,
	event *ingestionv1.UsageEvent,
	receivedAt time.Time,
) (ProcessResult, error) {
	if event == nil {
		return ProcessResult{}, shared.BadRequest("invalid_request", "usage event payload is required")
	}
	if receivedAt.IsZero() {
		receivedAt = time.Now().UTC()
	}
	return s.process(ctx, eventID, orgID, event, receivedAt)
}

func (s *Service) process(
	ctx context.Context,
	eventID, orgID string,
	event *ingestionv1.UsageEvent,
	receivedAt time.Time,
) (ProcessResult, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return ProcessResult{}, err
	}
	eventUUID, err := parseUUID(eventID)
	if err != nil {
		return ProcessResult{}, err
	}

	idempotencyKey := strings.TrimSpace(event.GetIdempotencyKey())
	if idempotencyKey == "" {
		return ProcessResult{}, shared.BadRequest("invalid_request", "idempotency_key is required")
	}

	if _, err := s.queries.CheckIdempotency(ctx, sqlcgen.CheckIdempotencyParams{
		OrgID:          orgUUID,
		IdempotencyKey: idempotencyKey,
	}); err == nil {
		return ProcessResult{Status: statusDuplicate, CostUSD: "0"}, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return ProcessResult{}, shared.Internal("check idempotency")
	}

	occurredAt := time.Now().UTC()
	if event.GetOccurredAt() != nil {
		occurredAt = event.GetOccurredAt().AsTime().UTC()
	}

	quote, err := s.pricing.CalculateCost(
		ctx,
		event.GetProvider(),
		event.GetModel(),
		event.GetInputTokens(),
		event.GetOutputTokens(),
		event.GetCacheReadTokens(),
		event.GetCacheWriteTokens(),
		occurredAt,
	)
	if err != nil {
		return ProcessResult{}, shared.Internal("calculate cost")
	}

	row := chinfra.UsageEventRow{
		ID:               uuid.UUID(eventUUID.Bytes),
		OrgID:            uuid.UUID(orgUUID.Bytes),
		IdempotencyKey:   idempotencyKey,
		Provider:         strings.ToLower(strings.TrimSpace(event.GetProvider())),
		Model:            strings.TrimSpace(event.GetModel()),
		InputTokens:      event.GetInputTokens(),
		OutputTokens:     event.GetOutputTokens(),
		CacheReadTokens:  event.GetCacheReadTokens(),
		CacheWriteTokens: event.GetCacheWriteTokens(),
		CostUSD:          quote.CostUSD,
		IsPriced:         quote.IsPriced,
		OccurredAt:       occurredAt,
		ReceivedAt:       receivedAt.UTC(),
		Metadata:         map[string]string{},
	}
	if meta := event.GetMetadata(); meta != nil {
		row.TeamID = optionalUUID(meta.GetTeamId())
		row.UserID = optionalUUID(meta.GetUserId())
		row.Project = meta.GetProject()
		row.Feature = meta.GetFeature()
		row.Customer = meta.GetCustomer()
		if row.TeamID != nil {
			row.Metadata["team_id"] = meta.GetTeamId()
		}
		if row.UserID != nil {
			row.Metadata["user_id"] = meta.GetUserId()
		}
	}

	if err := s.clickhouse.InsertUsageEvent(ctx, row); err != nil {
		return ProcessResult{}, shared.Internal("insert usage event")
	}

	_, err = s.queries.RecordIdempotency(ctx, sqlcgen.RecordIdempotencyParams{
		OrgID:          orgUUID,
		IdempotencyKey: idempotencyKey,
		EventID:        eventUUID,
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().UTC().Add(s.idempotencyTTL), Valid: true},
	})
	if err != nil {
		if isUniqueViolation(err) {
			return ProcessResult{Status: statusDuplicate, CostUSD: quote.CostUSD}, nil
		}
		return ProcessResult{}, shared.Internal("record idempotency")
	}

	s.log.Info("processed usage event",
		slog.String("event_id", eventID),
		slog.String("org_id", orgID),
		slog.String("cost_usd", quote.CostUSD),
		slog.Bool("is_priced", quote.IsPriced),
	)

	return ProcessResult{Status: statusProcessed, CostUSD: quote.CostUSD}, nil
}

func parseUUID(raw string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return pgtype.UUID{}, shared.BadRequest("invalid_uuid", "invalid id format")
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

func optionalUUID(raw string) *uuid.UUID {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return nil
	}
	return &parsed
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}