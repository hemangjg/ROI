package di

import (
	"context"
	"log/slog"

	ingestapp "github.com/ai-finops/ai-finops/services/ingestion-service/internal/application/ingestion"
	svcconfig "github.com/ai-finops/ai-finops/services/ingestion-service/internal/config"
	"github.com/ai-finops/ai-finops/services/ingestion-service/internal/app"
	"github.com/ai-finops/ai-finops/services/ingestion-service/internal/infrastructure/postgres"
	redinfra "github.com/ai-finops/ai-finops/services/ingestion-service/internal/infrastructure/redis"
	grpcapi "github.com/ai-finops/ai-finops/services/ingestion-service/internal/interfaces/grpc"
	httpapi "github.com/ai-finops/ai-finops/services/ingestion-service/internal/interfaces/http"
	"github.com/ai-finops/ai-finops/packages/config"
	"github.com/ai-finops/ai-finops/packages/logger"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "ingestion-service"

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

func provideRateLimiter(cfg svcconfig.Config, log *slog.Logger) (*redinfra.RateLimiter, error) {
	return redinfra.NewRateLimiter(cfg.RedisURL, cfg.RateLimitPerMinute, log)
}

func provideIngestionService(
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
	limiter *redinfra.RateLimiter,
	log *slog.Logger,
) *ingestapp.Service {
	return ingestapp.NewService(pool, queries, limiter, log)
}

func provideIngestionHandler(svc *ingestapp.Service) *httpapi.IngestionHandler {
	return httpapi.NewIngestionHandler(svc)
}

func provideGRPCServer(cfg svcconfig.Config, svc *ingestapp.Service) (*grpcapi.Server, error) {
	return grpcapi.NewServer(cfg.GRPCPort, grpcapi.NewIngestionServer(svc))
}

func provideRouter(
	cfg svcconfig.Config,
	log *slog.Logger,
	queries *sqlcgen.Queries,
	ingestionHandler *httpapi.IngestionHandler,
) chi.Router {
	return httpapi.NewRouter(cfg.ServiceName, log, queries, ingestionHandler)
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
	limiter *redinfra.RateLimiter,
) *app.App {
	return &app.App{
		Config:     cfg,
		Log:        log,
		Router:     router,
		Pool:       pool,
		Queries:    queries,
		GRPCServer: grpcServer,
		RateLimiter: limiter,
	}
}