package grpcapi

import (
	"context"
	"time"

	ingestapp "github.com/ai-finops/ai-finops/services/ingestion-service/internal/application/ingestion"
	ingestionv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/ingestion/v1"
)

// IngestionServer implements ingestion.v1.IngestionService.
type IngestionServer struct {
	ingestionv1.UnimplementedIngestionServiceServer
	svc *ingestapp.Service
}

// NewIngestionServer constructs the gRPC ingestion service.
func NewIngestionServer(svc *ingestapp.Service) *IngestionServer {
	return &IngestionServer{svc: svc}
}

// IngestEvent queues a usage event for an explicit organization.
func (s *IngestionServer) IngestEvent(ctx context.Context, req *ingestionv1.IngestEventRequest) (*ingestionv1.IngestEventResponse, error) {
	event, err := protoToUsageEvent(req.GetEvent())
	if err != nil {
		return nil, toGRPCError(err)
	}

	result, err := s.svc.IngestEventForOrg(ctx, req.GetOrgId(), event)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestionv1.IngestEventResponse{
		EventId: result.EventID,
		Status:  result.Status,
	}, nil
}

func protoToUsageEvent(event *ingestionv1.UsageEvent) (ingestapp.UsageEvent, error) {
	if event == nil {
		return ingestapp.UsageEvent{}, validationErr("event is required")
	}

	occurredAt := time.Now().UTC()
	if event.GetOccurredAt() != nil {
		occurredAt = event.GetOccurredAt().AsTime().UTC()
	}

	usageEvent := ingestapp.UsageEvent{
		IdempotencyKey:   event.GetIdempotencyKey(),
		Provider:         event.GetProvider(),
		Model:            event.GetModel(),
		InputTokens:      event.GetInputTokens(),
		OutputTokens:     event.GetOutputTokens(),
		CacheReadTokens:  event.GetCacheReadTokens(),
		CacheWriteTokens: event.GetCacheWriteTokens(),
		OccurredAt:       occurredAt,
	}

	if meta := event.GetMetadata(); meta != nil {
		usageEvent.Metadata = &ingestapp.UsageEventMetadata{
			TeamID:   meta.GetTeamId(),
			UserID:   meta.GetUserId(),
			Project:  meta.GetProject(),
			Feature:  meta.GetFeature(),
			Customer: meta.GetCustomer(),
		}
	}

	return usageEvent, nil
}