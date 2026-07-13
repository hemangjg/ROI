package domain_test

import (
	"testing"

	"github.com/ai-finops/ai-finops/services/pricing-service/internal/domain"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func TestCalculateTokenCost(t *testing.T) {
	inputPrice := mustNumeric(t, "2.50")
	outputPrice := mustNumeric(t, "10.00")
	cacheReadPrice := mustNumeric(t, "1.25")

	price := sqlcgen.ModelPricing{
		InputPricePer1m:     inputPrice,
		OutputPricePer1m:    outputPrice,
		CacheReadPricePer1m: cacheReadPrice,
	}

	cost, err := domain.CalculateTokenCost(1_000, 500, 200, 0, price)
	if err != nil {
		t.Fatalf("calculate cost: %v", err)
	}

	expected := decimal.RequireFromString("0.00775000")
	if !cost.Equal(expected) {
		t.Fatalf("cost = %s, want %s", cost.StringFixed(8), expected.StringFixed(8))
	}
}

func mustNumeric(t *testing.T, raw string) pgtype.Numeric {
	t.Helper()
	value, err := domain.DecimalToNumeric(decimal.RequireFromString(raw))
	if err != nil {
		t.Fatalf("decimal to numeric: %v", err)
	}
	return value
}