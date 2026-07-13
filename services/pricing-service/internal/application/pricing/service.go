package pricing

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/ai-finops/ai-finops/packages/shared"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/domain"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/plugin"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CostResult is returned by cost calculation.
type CostResult struct {
	CostUSD  string
	IsPriced bool
	ModelID  string
}

// Service implements pricing catalog sync and cost calculation.
type Service struct {
	pool     *pgxpool.Pool
	queries  *sqlcgen.Queries
	registry *plugin.Registry
	log      *slog.Logger
}

// NewService constructs a pricing application service.
func NewService(pool *pgxpool.Pool, queries *sqlcgen.Queries, registry *plugin.Registry, log *slog.Logger) *Service {
	return &Service{
		pool:     pool,
		queries:  queries,
		registry: registry,
		log:      log,
	}
}

// CalculateCost looks up versioned pricing and computes USD cost for a usage event.
func (s *Service) CalculateCost(
	ctx context.Context,
	provider, model string,
	inputTokens, outputTokens, cacheReadTokens, cacheWriteTokens uint32,
	occurredAt time.Time,
) (CostResult, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	model = strings.TrimSpace(model)
	if provider == "" || model == "" {
		return CostResult{}, shared.BadRequest("invalid_request", "provider and model are required")
	}
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	price, err := s.queries.GetPriceAtTime(ctx, sqlcgen.GetPriceAtTimeParams{
		Slug:          provider,
		Slug_2:        model,
		EffectiveFrom: pgtype.Timestamptz{Time: occurredAt.UTC(), Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CostResult{CostUSD: "0", IsPriced: false}, nil
		}
		return CostResult{}, shared.Internal("lookup model pricing")
	}

	cost, err := domain.CalculateTokenCost(inputTokens, outputTokens, cacheReadTokens, cacheWriteTokens, price)
	if err != nil {
		return CostResult{}, shared.Internal("calculate token cost")
	}

	return CostResult{
		CostUSD:  cost.StringFixed(8),
		IsPriced: true,
		ModelID:  pgUUIDString(price.ModelID),
	}, nil
}

// SyncPricingCatalog upserts providers, models, and versioned pricing from all plugins.
func (s *Service) SyncPricingCatalog(ctx context.Context) (uint32, error) {
	now := time.Now().UTC()
	var updated uint32

	for _, providerPlugin := range s.registry.All() {
		entries, err := providerPlugin.FetchPricing(ctx)
		if err != nil {
			return updated, shared.Internal("fetch pricing for " + providerPlugin.Name())
		}

		for _, entry := range entries {
			if err := s.upsertCatalogEntry(ctx, entry, now); err != nil {
				return updated, err
			}
			updated++
		}
	}

	return updated, nil
}

func (s *Service) upsertCatalogEntry(ctx context.Context, entry domain.CatalogEntry, effectiveFrom time.Time) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return shared.Internal("begin transaction")
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	provider, err := qtx.UpsertProvider(ctx, sqlcgen.UpsertProviderParams{
		Slug: entry.Provider,
		Name: providerDisplayName(entry.Provider),
	})
	if err != nil {
		return shared.Internal("upsert provider")
	}

	model, err := qtx.UpsertModel(ctx, sqlcgen.UpsertModelParams{
		ProviderID: provider.ID,
		Slug:       entry.Model,
		Name:       entry.DisplayName,
	})
	if err != nil {
		return shared.Internal("upsert model")
	}

	if err := qtx.CloseOpenModelPricing(ctx, sqlcgen.CloseOpenModelPricingParams{
		ModelID:     model.ID,
		EffectiveTo: pgtype.Timestamptz{Time: effectiveFrom, Valid: true},
	}); err != nil {
		return shared.Internal("close open model pricing")
	}

	inputPrice, err := domain.DecimalToNumeric(entry.InputPricePer1M)
	if err != nil {
		return shared.Internal("encode input price")
	}
	outputPrice, err := domain.DecimalToNumeric(entry.OutputPricePer1M)
	if err != nil {
		return shared.Internal("encode output price")
	}
	cacheReadPrice, err := domain.DecimalToNumeric(entry.CacheReadPricePer1M)
	if err != nil {
		return shared.Internal("encode cache read price")
	}
	cacheWritePrice, err := domain.DecimalToNumeric(entry.CacheWritePricePer1M)
	if err != nil {
		return shared.Internal("encode cache write price")
	}

	_, err = qtx.InsertModelPricing(ctx, sqlcgen.InsertModelPricingParams{
		ModelID:              model.ID,
		InputPricePer1m:      inputPrice,
		OutputPricePer1m:     outputPrice,
		CacheReadPricePer1m:  cacheReadPrice,
		CacheWritePricePer1m: cacheWritePrice,
		EffectiveFrom:        pgtype.Timestamptz{Time: effectiveFrom, Valid: true},
		EffectiveTo:          pgtype.Timestamptz{},
	})
	if err != nil {
		return shared.Internal("insert model pricing")
	}

	if err := tx.Commit(ctx); err != nil {
		return shared.Internal("commit transaction")
	}
	return nil
}

func providerDisplayName(slug string) string {
	switch slug {
	case "openai":
		return "OpenAI"
	case "anthropic":
		return "Anthropic"
	default:
		return strings.ToUpper(slug[:1]) + slug[1:]
	}
}

func pgUUIDString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	u, err := uuid.FromBytes(id.Bytes[:])
	if err != nil {
		return ""
	}
	return u.String()
}