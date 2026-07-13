package anthropic

import (
	"context"

	"github.com/ai-finops/ai-finops/services/pricing-service/internal/domain"
	"github.com/shopspring/decimal"
)

// Plugin provides Anthropic catalog pricing for Phase 1 models.
type Plugin struct{}

// NewPlugin constructs the Anthropic pricing plugin.
func NewPlugin() *Plugin {
	return &Plugin{}
}

// Name returns the provider slug.
func (p *Plugin) Name() string {
	return "anthropic"
}

// FetchPricing returns the embedded Anthropic pricing catalog.
func (p *Plugin) FetchPricing(_ context.Context) ([]domain.CatalogEntry, error) {
	return []domain.CatalogEntry{
		entry("claude-3-5-sonnet-20241022", "Claude 3.5 Sonnet", "3.00", "15.00", "0.30", "3.75"),
		entry("claude-3-5-haiku-20241022", "Claude 3.5 Haiku", "0.80", "4.00", "0.08", "1.00"),
		entry("claude-3-opus-20240229", "Claude 3 Opus", "15.00", "75.00", "0", "0"),
	}, nil
}

func entry(model, displayName, input, output, cacheRead, cacheWrite string) domain.CatalogEntry {
	return domain.CatalogEntry{
		Provider:             "anthropic",
		Model:                model,
		DisplayName:          displayName,
		InputPricePer1M:      mustDecimal(input),
		OutputPricePer1M:     mustDecimal(output),
		CacheReadPricePer1M:  mustDecimal(cacheRead),
		CacheWritePricePer1M: mustDecimal(cacheWrite),
	}
}

func mustDecimal(raw string) decimal.Decimal {
	value, err := decimal.NewFromString(raw)
	if err != nil {
		panic("invalid catalog decimal: " + raw)
	}
	return value
}