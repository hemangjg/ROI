package di

import (
	"context"
	"log/slog"

	analyticsapp "github.com/ai-finops/ai-finops/services/analytics-service/internal/application/analytics"
	svcconfig "github.com/ai-finops/ai-finops/services/analytics-service/internal/config"
	"github.com/ai-finops/ai-finops/services/analytics-service/internal/app"
	chinfra "github.com/ai-finops/ai-finops/services/analytics-service/internal/infrastructure/clickhouse"
	"github.com/ai-finops/ai-finops/services/analytics-service/internal/infrastructure/postgres"
	grpcapi "github.com/ai-finops/ai-finops/services/analytics-service/internal/interfaces/grpc"
	httpapi "github.com/ai-finops/ai-finops/services/analytics-service/internal/interfaces/http"
	"github.com/ai-finops/ai-finops/packages/config"
	"github.com/ai-finops/ai-finops/packages/logger"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "analytics-service"

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

func provideClickHouseClient(cfg svcconfig.Config) (*chinfra.Client, error) {
	return chinfra.NewClient(cfg.ClickHouseURL)
}

func provideAnalyticsService(
	clickhouse *chinfra.Client,
	queries *sqlcgen.Queries,
	log *slog.Logger,
) *analyticsapp.Service {
	return analyticsapp.NewService(clickhouse, queries, log)
}

func provideAnalyticsHandler(svc *analyticsapp.Service) *httpapi.AnalyticsHandler {
	return httpapi.NewAnalyticsHandler(svc)
}

func provideGRPCServer(cfg svcconfig.Config, svc *analyticsapp.Service) (*grpcapi.Server, error) {
	return grpcapi.NewServer(cfg.GRPCPort, grpcapi.NewAnalyticsServer(svc))
}

func provideRouter(
	cfg svcconfig.Config,
	log *slog.Logger,
	queries *sqlcgen.Queries,
	analytics *httpapi.AnalyticsHandler,
) chi.Router {
	return httpapi.NewRouter(cfg.ServiceName, log, queries, analytics)
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
	clickhouse *chinfra.Client,
) *app.App {
	return &app.App{
		Config:     cfg,
		Log:        log,
		Router:     router,
		Pool:       pool,
		Queries:    queries,
		GRPCServer: grpcServer,
		ClickHouse: clickhouse,
	}
}