package anthropic_test

import (
	"context"
	"testing"

	"github.com/ai-finops/ai-finops/services/pricing-service/internal/plugin/anthropic"
)

func TestAnthropicPluginCatalog(t *testing.T) {
	p := anthropic.NewPlugin()
	if p.Name() != "anthropic" {
		t.Fatalf("name = %q", p.Name())
	}
	entries, err := p.FetchPricing(context.Background())
	if err != nil {
		t.Fatalf("FetchPricing: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected non-empty anthropic catalog")
	}
}
