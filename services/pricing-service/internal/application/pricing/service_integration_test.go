package pricing_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ai-finops/ai-finops/packages/logger"
	pricingapp "github.com/ai-finops/ai-finops/services/pricing-service/internal/application/pricing"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/infrastructure/postgres"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/plugin"
	"github.com/shopspring/decimal"
)

func TestPricingServiceIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, queries, err := postgres.Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer pool.Close()

	svc := pricingapp.NewService(pool, queries, plugin.NewRegistry(), logger.New("pricing-service-test", "error"))

	updated, err := svc.SyncPricingCatalog(ctx)
	if err != nil {
		t.Fatalf("sync catalog: %v", err)
	}
	if updated < 8 {
		t.Fatalf("expected at least 8 models updated, got %d", updated)
	}

	occurredAt := time.Now().UTC()
	result, err := svc.CalculateCost(ctx, "openai", "gpt-4o-mini", 1_000_000, 0, 0, 0, occurredAt)
	if err != nil {
		t.Fatalf("calculate cost: %v", err)
	}
	if !result.IsPriced || result.ModelID == "" {
		t.Fatalf("expected priced result, got %+v", result)
	}

	expected := decimal.RequireFromString("0.15000000")
	actual, err := decimal.NewFromString(result.CostUSD)
	if err != nil {
		t.Fatalf("parse cost: %v", err)
	}
	if !actual.Equal(expected) {
		t.Fatalf("cost = %s, want %s", result.CostUSD, expected.StringFixed(8))
	}

	unknown, err := svc.CalculateCost(ctx, "openai", "unknown-model", 100, 0, 0, 0, occurredAt)
	if err != nil {
		t.Fatalf("calculate unknown model: %v", err)
	}
	if unknown.IsPriced {
		t.Fatalf("expected unpriced unknown model, got %+v", unknown)
	}
}