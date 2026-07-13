package events_test

import (
	"testing"

	"github.com/ai-finops/ai-finops/packages/events"
)

func TestEnsureTopic_Validation(t *testing.T) {
	if err := events.EnsureTopic("", events.TopicUsageEventsRaw, 1); err == nil {
		t.Fatal("expected error for empty brokers")
	}
	if err := events.EnsureTopic("localhost:9092", "", 1); err == nil {
		t.Fatal("expected error for empty topic")
	}
}
