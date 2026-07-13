package openai_test

import (
	"context"
	"testing"

	"github.com/ai-finops/ai-finops/services/pricing-service/internal/plugin/openai"
)

func TestOpenAIPluginCatalog(t *testing.T) {
	p := openai.NewPlugin()
	if p.Name() != "openai" {
		t.Fatalf("name = %q", p.Name())
	}
	entries, err := p.FetchPricing(context.Background())
	if err != nil {
		t.Fatalf("FetchPricing: %v", err)
	}
	if len(entries) < 3 {
		t.Fatalf("expected catalog entries, got %d", len(entries))
	}
	found := false
	for _, e := range entries {
		if e.Model == "gpt-4o" {
			found = true
			if e.InputPricePer1M.IsZero() || e.OutputPricePer1M.IsZero() {
				t.Fatal("gpt-4o prices must be non-zero")
			}
		}
	}
	if !found {
		t.Fatal("expected gpt-4o in catalog")
	}
}
