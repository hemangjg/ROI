package relay

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPayloadToProto(t *testing.T) {
	payload := map[string]any{
		"event_id":    "11111111-1111-1111-1111-111111111111",
		"org_id":      "22222222-2222-2222-2222-222222222222",
		"received_at": time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC),
		"event": map[string]any{
			"idempotency_key":    "key-1",
			"provider":           "openai",
			"model":              "gpt-4o",
			"input_tokens":       100,
			"output_tokens":      50,
			"cache_read_tokens":  0,
			"cache_write_tokens": 0,
			"occurred_at":        time.Date(2026, 7, 7, 11, 59, 0, 0, time.UTC),
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	message, err := payloadToProto("usage.event.created", body)
	if err != nil {
		t.Fatalf("convert payload: %v", err)
	}
	if message.GetEventId() == "" || message.GetOrgId() == "" {
		t.Fatalf("expected event and org ids, got %+v", message)
	}
	if message.GetEvent().GetProvider() != "openai" {
		t.Fatalf("provider = %q", message.GetEvent().GetProvider())
	}
}