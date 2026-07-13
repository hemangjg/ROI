package openai

import (
	"context"

	"github.com/ai-finops/ai-finops/services/pricing-service/internal/domain"
	"github.com/shopspring/decimal"
)

// Plugin provides OpenAI catalog pricing for Phase 1 models.
type Plugin struct{}

// NewPlugin constructs the OpenAI pricing plugin.
func NewPlugin() *Plugin {
	return &Plugin{}
}

// Name returns the provider slug.
func (p *Plugin) Name() string {
	return "openai"
}

// FetchPricing returns the embedded OpenAI pricing catalog.
func (p *Plugin) FetchPricing(_ context.Context) ([]domain.CatalogEntry, error) {
	return []domain.CatalogEntry{
		entry("gpt-4o", "GPT-4o", "2.50", "10.00", "1.25", "0"),
		entry("gpt-4o-mini", "GPT-4o Mini", "0.15", "0.60", "0.075", "0"),
		entry("gpt-4-turbo", "GPT-4 Turbo", "10.00", "30.00", "0", "0"),
		entry("o1", "o1", "15.00", "60.00", "0", "0"),
		entry("o1-mini", "o1 Mini", "3.00", "12.00", "0", "0"),
	}, nil
}

func entry(model, displayName, input, output, cacheRead, cacheWrite string) domain.CatalogEntry {
	return domain.CatalogEntry{
		Provider:             "openai",
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