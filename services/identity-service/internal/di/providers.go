package di

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	authapp "github.com/ai-finops/ai-finops/services/identity-service/internal/application/auth"
	svcconfig "github.com/ai-finops/ai-finops/services/identity-service/internal/config"
	"github.com/ai-finops/ai-finops/services/identity-service/internal/app"
	"github.com/ai-finops/ai-finops/services/identity-service/internal/infrastructure/jwt"
	"github.com/ai-finops/ai-finops/services/identity-service/internal/infrastructure/postgres"
	grpcapi "github.com/ai-finops/ai-finops/services/identity-service/internal/interfaces/grpc"
	httpapi "github.com/ai-finops/ai-finops/services/identity-service/internal/interfaces/http"
	"github.com/ai-finops/ai-finops/packages/config"
	"github.com/ai-finops/ai-finops/packages/logger"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "identity-service"

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

func provideJWTKeys(cfg svcconfig.Config) (*jwt.Keys, error) {
	return jwt.Load(cfg.JWTPrivateKeyPath, cfg.JWTPublicKeyPath)
}

func provideAuthService(
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
	keys *jwt.Keys,
	cfg svcconfig.Config,
	log *slog.Logger,
) (*authapp.Service, error) {
	accessTTL, err := time.ParseDuration(cfg.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("parse access_token_ttl: %w", err)
	}
	refreshTTL, err := time.ParseDuration(cfg.RefreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("parse refresh_token_ttl: %w", err)
	}

	return authapp.NewService(pool, queries, keys.Private, keys.Public, accessTTL, refreshTTL, log), nil
}

func provideAuthHandler(svc *authapp.Service) *httpapi.AuthHandler {
	return httpapi.NewAuthHandler(svc)
}

func provideGRPCServer(cfg svcconfig.Config, svc *authapp.Service) (*grpcapi.Server, error) {
	return grpcapi.NewServer(
		cfg.GRPCPort,
		grpcapi.NewAuthServer(svc),
		grpcapi.NewExtAuthzServer(svc),
	)
}

func provideRouter(
	cfg svcconfig.Config,
	log *slog.Logger,
	queries *sqlcgen.Queries,
	authHandler *httpapi.AuthHandler,
) chi.Router {
	return httpapi.NewRouter(cfg.ServiceName, log, queries, authHandler)
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