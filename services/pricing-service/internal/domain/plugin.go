package domain

import (
	"context"

	"github.com/shopspring/decimal"
)

// CatalogEntry is a provider model price row used during catalog sync.
type CatalogEntry struct {
	Provider             string
	Model                string
	DisplayName          string
	InputPricePer1M      decimal.Decimal
	OutputPricePer1M     decimal.Decimal
	CacheReadPricePer1M  decimal.Decimal
	CacheWritePricePer1M decimal.Decimal
}

// ProviderPlugin supplies static or fetched pricing for a vendor.
type ProviderPlugin interface {
	Name() string
	FetchPricing(ctx context.Context) ([]CatalogEntry, error)
}