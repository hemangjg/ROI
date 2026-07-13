package di

import (
	"context"
	"log/slog"

	mgmtapp "github.com/ai-finops/ai-finops/services/management-service/internal/application/management"
	svcconfig "github.com/ai-finops/ai-finops/services/management-service/internal/config"
	"github.com/ai-finops/ai-finops/services/management-service/internal/app"
	"github.com/ai-finops/ai-finops/services/management-service/internal/infrastructure/postgres"
	grpcapi "github.com/ai-finops/ai-finops/services/management-service/internal/interfaces/grpc"
	httpapi "github.com/ai-finops/ai-finops/services/management-service/internal/interfaces/http"
	"github.com/ai-finops/ai-finops/packages/config"
	"github.com/ai-finops/ai-finops/packages/logger"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "management-service"

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

func provideManagementService(
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
	log *slog.Logger,
) *mgmtapp.Service {
	return mgmtapp.NewService(pool, queries, log)
}

func provideManagementHandler(svc *mgmtapp.Service) *httpapi.ManagementHandler {
	return httpapi.NewManagementHandler(svc)
}

func provideGRPCServer(cfg svcconfig.Config, svc *mgmtapp.Service) (*grpcapi.Server, error) {
	return grpcapi.NewServer(cfg.GRPCPort, grpcapi.NewManagementServer(svc))
}

func provideRouter(
	cfg svcconfig.Config,
	log *slog.Logger,
	queries *sqlcgen.Queries,
	mgmtHandler *httpapi.ManagementHandler,
) chi.Router {
	return httpapi.NewRouter(cfg.ServiceName, log, queries, mgmtHandler)
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