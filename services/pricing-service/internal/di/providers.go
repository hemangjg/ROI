package di

import (
	"context"
	"log/slog"

	pricingapp "github.com/ai-finops/ai-finops/services/pricing-service/internal/application/pricing"
	svcconfig "github.com/ai-finops/ai-finops/services/pricing-service/internal/config"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/app"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/infrastructure/postgres"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/plugin"
	grpcapi "github.com/ai-finops/ai-finops/services/pricing-service/internal/interfaces/grpc"
	httpapi "github.com/ai-finops/ai-finops/services/pricing-service/internal/interfaces/http"
	"github.com/ai-finops/ai-finops/packages/config"
	"github.com/ai-finops/ai-finops/packages/logger"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "pricing-service"

func provideConfig() (svcconfig.Config, error) {
	return config.Load[svcconfig.Config](serviceName, "config/defaults.yaml", "config/config.local.yaml")
}

func provideLogger(cfg svcconfig.Config) *slog.Logger {
	return logger.New(cfg.ServiceName, cfg.LogLevel)
}

func providePool(cfg svcconfig.Config) (*pgxpool.Pool, error) {
	pool, _, err := postgres.Open(context.Background(), cfg.DatabaseURL)
	return pool, err
}

func provideQueries(pool *pgxpool.Pool) *sqlcgen.Queries {
	return sqlcgen.New(pool)
}

func providePluginRegistry() *plugin.Registry {
	return plugin.NewRegistry()
}

func providePricingService(
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
	registry *plugin.Registry,
	log *slog.Logger,
) *pricingapp.Service {
	return pricingapp.NewService(pool, queries, registry, log)
}

func provideGRPCServer(cfg svcconfig.Config, svc *pricingapp.Service) (*grpcapi.Server, error) {
	return grpcapi.NewServer(cfg.GRPCPort, grpcapi.NewPricingServer(svc))
}

func provideRouter(cfg svcconfig.Config, log *slog.Logger, queries *sqlcgen.Queries) chi.Router {
	return httpapi.NewRouter(cfg.ServiceName, log, queries)
}

func provideTelemetry(cfg svcconfig.Config, log *slog.Logger) (func(context.Context) error, error) {
	return shared.InitTelemetry(context.Background(), cfg.ServiceName, cfg.Environment, cfg.OTLPEndpoint, log)
}

func provideApp(
	cfg svcconfig.Config,
	log *slog.Logger,
	router chi.Router,
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
	grpcServer *grpcapi.Server,
) *app.App {
	return &app.App{
		Config:     cfg,
		Log:        log,
		Router:     router,
		Pool:       pool,
		Queries:    queries,
		GRPCServer: grpcServer,
	}
}