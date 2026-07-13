package relay

import (
	"encoding/json"
	"fmt"
	"time"

	ingestionv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/ingestion/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type outboxPayload struct {
	EventID    string          `json:"event_id"`
	OrgID      string          `json:"org_id"`
	Event      outboxUsageEvent `json:"event"`
	ReceivedAt time.Time       `json:"received_at"`
}

type outboxUsageEvent struct {
	IdempotencyKey   string                 `json:"idempotency_key"`
	Provider         string                 `json:"provider"`
	Model            string                 `json:"model"`
	InputTokens      uint32                 `json:"input_tokens"`
	OutputTokens     uint32                 `json:"output_tokens"`
	CacheReadTokens  uint32                 `json:"cache_read_tokens"`
	CacheWriteTokens uint32                 `json:"cache_write_tokens"`
	OccurredAt       time.Time              `json:"occurred_at"`
	Metadata         *outboxEventMetadata   `json:"metadata,omitempty"`
}

type outboxEventMetadata struct {
	TeamID   string `json:"team_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	Project  string `json:"project,omitempty"`
	Feature  string `json:"feature,omitempty"`
	Customer string `json:"customer,omitempty"`
}

func payloadToProto(eventType string, payload []byte) (*ingestionv1.UsageEventCreated, error) {
	if eventType != "usage.event.created" {
		return nil, fmt.Errorf("unsupported event type %q", eventType)
	}

	var raw outboxPayload
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, fmt.Errorf("unmarshal outbox payload: %w", err)
	}

	usageEvent := &ingestionv1.UsageEvent{
		IdempotencyKey:   raw.Event.IdempotencyKey,
		Provider:         raw.Event.Provider,
		Model:            raw.Event.Model,
		InputTokens:      raw.Event.InputTokens,
		OutputTokens:     raw.Event.OutputTokens,
		CacheReadTokens:  raw.Event.CacheReadTokens,
		CacheWriteTokens: raw.Event.CacheWriteTokens,
		OccurredAt:       timestamppb.New(raw.Event.OccurredAt.UTC()),
	}
	if raw.Event.Metadata != nil {
		usageEvent.Metadata = &ingestionv1.UsageEventMetadata{
			TeamId:   raw.Event.Metadata.TeamID,
			UserId:   raw.Event.Metadata.UserID,
			Project:  raw.Event.Metadata.Project,
			Feature:  raw.Event.Metadata.Feature,
			Customer: raw.Event.Metadata.Customer,
		}
	}

	return &ingestionv1.UsageEventCreated{
		EventId:    raw.EventID,
		OrgId:      raw.OrgID,
		Event:      usageEvent,
		ReceivedAt: timestamppb.New(raw.ReceivedAt.UTC()),
	}, nil
}