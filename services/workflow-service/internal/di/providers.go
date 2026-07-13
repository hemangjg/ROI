package di

import (
	"context"
	"log/slog"
	"time"

	workflowapp "github.com/ai-finops/ai-finops/services/workflow-service/internal/application/workflow"
	svcconfig "github.com/ai-finops/ai-finops/services/workflow-service/internal/config"
	"github.com/ai-finops/ai-finops/services/workflow-service/internal/app"
	chinfra "github.com/ai-finops/ai-finops/services/workflow-service/internal/infrastructure/clickhouse"
	kafkainfra "github.com/ai-finops/ai-finops/services/workflow-service/internal/infrastructure/kafka"
	"github.com/ai-finops/ai-finops/services/workflow-service/internal/infrastructure/postgres"
	pricinginfra "github.com/ai-finops/ai-finops/services/workflow-service/internal/infrastructure/pricing"
	grpcapi "github.com/ai-finops/ai-finops/services/workflow-service/internal/interfaces/grpc"
	httpapi "github.com/ai-finops/ai-finops/services/workflow-service/internal/interfaces/http"
	"github.com/ai-finops/ai-finops/packages/config"
	"github.com/ai-finops/ai-finops/packages/logger"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "workflow-service"

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

func providePricingClient(cfg svcconfig.Config) (*pricinginfra.Client, error) {
	return pricinginfra.NewClient(cfg.PricingGRPCAddr)
}

func provideClickHouseClient(cfg svcconfig.Config) (*chinfra.Client, error) {
	return chinfra.NewClient(cfg.ClickHouseURL)
}

func provideWorkflowService(
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
	pricing *pricinginfra.Client,
	clickhouse *chinfra.Client,
	cfg svcconfig.Config,
	log *slog.Logger,
) *workflowapp.Service {
	ttl := time.Duration(cfg.IdempotencyTTLHours) * time.Hour
	return workflowapp.NewService(pool, queries, pricing, clickhouse, ttl, log)
}

func provideKafkaConsumer(cfg svcconfig.Config, svc *workflowapp.Service, log *slog.Logger) *kafkainfra.Consumer {
	return kafkainfra.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID, svc, log)
}

func provideGRPCServer(cfg svcconfig.Config, svc *workflowapp.Service) (*grpcapi.Server, error) {
	return grpcapi.NewServer(cfg.GRPCPort, grpcapi.NewWorkflowServer(svc))
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
	kafkaConsumer *kafkainfra.Consumer,
	pricing *pricinginfra.Client,
	clickhouse *chinfra.Client,
) *app.App {
	return &app.App{
		Config:        cfg,
		Log:           log,
		Router:        router,
		Pool:          pool,
		Queries:       queries,
		GRPCServer:    grpcServer,
		KafkaConsumer: kafkaConsumer,
		PricingClient: pricing,
		ClickHouse:    clickhouse,
	}
}