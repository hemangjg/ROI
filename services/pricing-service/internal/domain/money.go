package domain

import (
	"fmt"

	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

const tokensPerMillion = 1_000_000

// CalculateTokenCost computes USD cost from token counts and per-1M prices.
func CalculateTokenCost(
	inputTokens, outputTokens, cacheReadTokens, cacheWriteTokens uint32,
	price sqlcgen.ModelPricing,
) (decimal.Decimal, error) {
	inputPrice, err := numericToDecimal(price.InputPricePer1m)
	if err != nil {
		return decimal.Zero, fmt.Errorf("input price: %w", err)
	}
	outputPrice, err := numericToDecimal(price.OutputPricePer1m)
	if err != nil {
		return decimal.Zero, fmt.Errorf("output price: %w", err)
	}
	cacheReadPrice, err := numericToDecimal(price.CacheReadPricePer1m)
	if err != nil {
		return decimal.Zero, fmt.Errorf("cache read price: %w", err)
	}
	cacheWritePrice, err := numericToDecimal(price.CacheWritePricePer1m)
	if err != nil {
		return decimal.Zero, fmt.Errorf("cache write price: %w", err)
	}

	million := decimal.NewFromInt(tokensPerMillion)
	cost := decimal.Zero
	cost = cost.Add(inputPrice.Mul(decimal.NewFromInt(int64(inputTokens))).Div(million))
	cost = cost.Add(outputPrice.Mul(decimal.NewFromInt(int64(outputTokens))).Div(million))
	cost = cost.Add(cacheReadPrice.Mul(decimal.NewFromInt(int64(cacheReadTokens))).Div(million))
	cost = cost.Add(cacheWritePrice.Mul(decimal.NewFromInt(int64(cacheWriteTokens))).Div(million))
	return cost, nil
}

// DecimalToNumeric converts a decimal value for PostgreSQL NUMERIC columns.
func DecimalToNumeric(value decimal.Decimal) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	if err := numeric.Scan(value.StringFixed(8)); err != nil {
		return pgtype.Numeric{}, err
	}
	return numeric, nil
}

func numericToDecimal(value pgtype.Numeric) (decimal.Decimal, error) {
	if !value.Valid {
		return decimal.Zero, nil
	}
	floatValue, err := value.Float64Value()
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromFloat(floatValue.Float64), nil
}