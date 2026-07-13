package events

import "time"

// Topic names for the usage event pipeline.
const (
	TopicUsageEventsRaw = "usage.events.raw"
	TopicUsageDLQ       = "usage.events.dlq"
)

// Envelope is a generic event wrapper skeleton.
type Envelope struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	OrgID     string         `json:"org_id"`
	Timestamp time.Time      `json:"timestamp"`
	Payload   map[string]any `json:"payload,omitempty"`
}