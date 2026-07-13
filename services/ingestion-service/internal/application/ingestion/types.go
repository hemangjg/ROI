package ingestion

import "time"

const (
	maxBatchSize       = 500
	maxIdempotencyLen  = 128
	outboxEventType    = "usage.event.created"
	statusQueued       = "queued"
)

// UsageEventMetadata is optional attribution metadata on a usage event.
type UsageEventMetadata struct {
	TeamID   string `json:"team_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	Project  string `json:"project,omitempty"`
	Feature  string `json:"feature,omitempty"`
	Customer string `json:"customer,omitempty"`
}

// UsageEvent is the REST request body for ingestion.
type UsageEvent struct {
	IdempotencyKey   string              `json:"idempotency_key" validate:"required,max=128"`
	Provider         string              `json:"provider" validate:"required"`
	Model            string              `json:"model" validate:"required"`
	InputTokens      uint32              `json:"input_tokens" validate:"gte=0"`
	OutputTokens     uint32              `json:"output_tokens" validate:"gte=0"`
	CacheReadTokens  uint32              `json:"cache_read_tokens"`
	CacheWriteTokens uint32              `json:"cache_write_tokens"`
	OccurredAt       time.Time           `json:"occurred_at" validate:"required"`
	Metadata         *UsageEventMetadata `json:"metadata,omitempty"`
}

// IngestResult is returned when an event is accepted.
type IngestResult struct {
	EventID string `json:"event_id"`
	Status  string `json:"status"`
}

// BatchIngestResult summarizes a batch ingest call.
type BatchIngestResult struct {
	Accepted int                 `json:"accepted"`
	Rejected int                 `json:"rejected"`
	Errors   []BatchIngestError  `json:"errors,omitempty"`
}

// BatchIngestError describes a rejected event in a batch.
type BatchIngestError struct {
	Index   int    `json:"index"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// EventStatus is returned by the event status endpoint.
type EventStatus struct {
	EventID string `json:"event_id"`
	Status  string `json:"status"`
	CostUSD string `json:"cost_usd,omitempty"`
}

// OutboxPayload is persisted in outbox_events.payload for the relay consumer.
type OutboxPayload struct {
	EventID    string     `json:"event_id"`
	OrgID      string     `json:"org_id"`
	Event      UsageEvent `json:"event"`
	ReceivedAt time.Time  `json:"received_at"`
}

// APIKeyIdentity is the authenticated ingestion caller.
type APIKeyIdentity struct {
	OrgID string
	KeyID string
}