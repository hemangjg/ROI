#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

declare -A SERVICE_PORTS=(
  [identity-service]=8080
  [ingestion-service]=8081
  [management-service]=8082
  [analytics-service]=8083
  [pricing-service]=8084
  [workflow-service]=8085
)

write_file() {
  local path="$1"
  mkdir -p "$(dirname "$path")"
  cat > "$path"
}

for service in "${!SERVICE_PORTS[@]}"; do
  port="${SERVICE_PORTS[$service]}"
  svc_dir="$ROOT/services/$service"
  mod="github.com/ai-finops/ai-finops/services/$service"

  echo "==> scaffolding $service (port $port)"

  write_file "$svc_dir/config/defaults.yaml" <<EOF
environment: local
log_level: info
http_port: $port
database_url: postgres://finops:finops@localhost:5432/finops?sslmode=disable
redis_url: redis://localhost:6379
otlp_endpoint: http://localhost:4317
EOF

  write_file "$svc_dir/go.mod" <<EOF
module $mod

go 1.25

require (
	github.com/ai-finops/ai-finops/db/atlas/sql v0.0.0
	github.com/ai-finops/ai-finops/packages/config v0.0.0
	github.com/ai-finops/ai-finops/packages/logger v0.0.0
	github.com/ai-finops/ai-finops/packages/shared v0.0.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/google/wire v0.6.0
	github.com/jackc/pgx/v5 v5.7.4
)

replace (
	github.com/ai-finops/ai-finops/db/atlas/sql => ../../db/atlas/sql
	github.com/ai-finops/ai-finops/packages/config => ../../packages/config
	github.com/ai-finops/ai-finops/packages/logger => ../../packages/logger
	github.com/ai-finops/ai-finops/packages/shared => ../../packages/shared
)
EOF

  write_file "$svc_dir/project.json" <<EOF
{
  "name": "$service",
  "\$schema": "../../node_modules/nx/schemas/project-schema.json",
  "projectType": "application",
  "sourceRoot": "services/$service",
  "tags": ["type:service", "scope:backend"],
  "targets": {
    "build": {
      "executor": "nx:run-commands",
      "options": {
        "command": "go build -o ../../dist/$service ./cmd/server",
        "cwd": "services/$service"
      }
    },
    "serve": {
      "executor": "nx:run-commands",
      "options": {
        "command": "go run ./cmd/server",
        "cwd": "services/$service"
      }
    },
    "test": {
      "executor": "nx:run-commands",
      "options": {
        "command": "go test ./... -count=1",
        "cwd": "services/$service"
      }
    },
    "lint": {
      "executor": "nx:run-commands",
      "options": {
        "command": "go vet ./...",
        "cwd": "services/$service"
      }
    }
  }
}
EOF

  write_file "$svc_dir/Dockerfile" <<EOF
FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.work ./
COPY packages ./packages
COPY db ./db
COPY services/$service ./services/$service

WORKDIR /src/services/$service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server /usr/local/bin/server
COPY services/$service/config/defaults.yaml /config/defaults.yaml
EXPOSE $port
ENTRYPOINT ["/usr/local/bin/server"]
EOF

  write_file "$svc_dir/internal/app/app.go" <<'EOF'
package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ai-finops/ai-finops/packages/config"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App holds runtime dependencies for the service shell.
type App struct {
	Config  config.Base
	Log     *slog.Logger
	Router  chi.Router
	Pool    *pgxpool.Pool
	Queries *sqlcgen.Queries
}

// HTTPServer returns the configured HTTP server.
func (a *App) HTTPServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", a.Config.HTTPPort),
		Handler:      a.Router,
		ReadTimeout:  15 * 1000000000,
		WriteTimeout: 15 * 1000000000,
		IdleTimeout:  60 * 1000000000,
	}
}
EOF

  # Fix app.go - I used nanoseconds incorrectly for timeouts. Let me fix in the script - use time package

  write_file "$svc_dir/internal/app/app.go" <<'EOF'
package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ai-finops/ai-finops/packages/config"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App holds runtime dependencies for the service shell.
type App struct {
	Config  config.Base
	Log     *slog.Logger
	Router  chi.Router
	Pool    *pgxpool.Pool
	Queries *sqlcgen.Queries
}

// HTTPServer returns the configured HTTP server.
func (a *App) HTTPServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", a.Config.HTTPPort),
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
EOF

  write_file "$svc_dir/internal/infrastructure/postgres/postgres.go" <<'EOF'
package postgres

import (
	"context"
	"fmt"

	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Open creates a PostgreSQL connection pool and SQLC queries handle.
func Open(ctx context.Context, databaseURL string) (*pgxpool.Pool, *sqlcgen.Queries, error) {
	if databaseURL == "" {
		return nil, nil, fmt.Errorf("database_url is required")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("connect postgres: %w", err)
	}

	return pool, sqlcgen.New(pool), nil
}
EOF

  write_file "$svc_dir/internal/interfaces/http/health.go" <<'EOF'
package httpapi

import (
	"encoding/json"
	"net/http"

	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
)

type healthResponse struct {
	Status string `json:"status"`
}

// Healthz reports process liveness.
func Healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}

// Readyz reports readiness, including database connectivity when configured.
func Readyz(queries *sqlcgen.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if queries == nil {
			writeJSON(w, http.StatusServiceUnavailable, healthResponse{Status: "database unavailable"})
			return
		}

		if _, err := queries.Ping(r.Context()); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, healthResponse{Status: "database unavailable"})
			return
		}

		writeJSON(w, http.StatusOK, healthResponse{Status: "ready"})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
EOF

  write_file "$svc_dir/internal/interfaces/http/router.go" <<'EOF'
package httpapi

import (
	"log/slog"
	"net/http"

	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/ai-finops/ai-finops/packages/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter wires the Phase 1 HTTP shell routes.
func NewRouter(serviceName string, log *slog.Logger, queries *sqlcgen.Queries) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(shared.OTelHTTP(serviceName))
	r.Use(middleware.Recoverer)
	r.Use(shared.RequestID)
	r.Use(shared.RequestLogger(log))

	r.Get("/healthz", Healthz)
	r.Get("/readyz", Readyz(queries))

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return r
}
EOF

  write_file "$svc_dir/internal/interfaces/http/health_test.go" <<'EOF'
package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpapi "github.com/ai-finops/ai-finops/services/SERVICE_NAME/internal/interfaces/http"
)

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	httpapi.Healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
EOF
  sed -i '' "s/SERVICE_NAME/$service/g" "$svc_dir/internal/interfaces/http/health_test.go" 2>/dev/null || \
    sed -i "s/SERVICE_NAME/$service/g" "$svc_dir/internal/interfaces/http/health_test.go"

  write_file "$svc_dir/internal/di/providers.go" <<EOF
package di

import (
	"context"
	"log/slog"

	"github.com/ai-finops/ai-finops/packages/config"
	"github.com/ai-finops/ai-finops/packages/logger"
	"$mod/internal/app"
	"$mod/internal/infrastructure/postgres"
	httpapi "$mod/internal/interfaces/http"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "$service"

func provideConfig() (config.Base, error) {
	return config.Load[config.Base](serviceName, "config/defaults.yaml", "config/config.local.yaml")
}

func provideLogger(cfg config.Base) *slog.Logger {
	return logger.New(cfg.ServiceName, cfg.LogLevel)
}

func providePool(cfg config.Base) (*pgxpool.Pool, error) {
	pool, _, err := postgres.Open(context.Background(), cfg.DatabaseURL)
	return pool, err
}

func provideQueries(pool *pgxpool.Pool) *sqlcgen.Queries {
	return sqlcgen.New(pool)
}

func provideRouter(cfg config.Base, log *slog.Logger, queries *sqlcgen.Queries) chi.Router {
	return httpapi.NewRouter(cfg.ServiceName, log, queries)
}

func provideTelemetry(cfg config.Base, log *slog.Logger) (func(context.Context) error, error) {
	return shared.InitTelemetry(context.Background(), cfg.ServiceName, cfg.Environment, cfg.OTLPEndpoint, log)
}

func provideApp(
	cfg config.Base,
	log *slog.Logger,
	router chi.Router,
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
) *app.App {
	return &app.App{
		Config:  cfg,
		Log:     log,
		Router:  router,
		Pool:    pool,
		Queries: queries,
	}
}
EOF

  write_file "$svc_dir/internal/di/wire.go" <<'EOF'
//go:build wireinject

package di

import (
	"github.com/google/wire"
	"SERVICE_MODULE/internal/app"
)

func Initialize() (*app.App, func(), error) {
	wire.Build(
		provideConfig,
		provideLogger,
		providePool,
		provideQueries,
		provideRouter,
		provideTelemetry,
		provideApp,
		newCleanup,
	)
	return nil, nil, nil
}
EOF
  sed -i '' "s|SERVICE_MODULE|$mod|g" "$svc_dir/internal/di/wire.go" 2>/dev/null || \
    sed -i "s|SERVICE_MODULE|$mod|g" "$svc_dir/internal/di/wire.go"

  write_file "$svc_dir/internal/di/cleanup.go" <<'EOF'
package di

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"SERVICE_MODULE/internal/app"
)

func newCleanup(
	app *app.App,
	pool *pgxpool.Pool,
	telemetryShutdown func(context.Context) error,
) func() {
	return func() {
		ctx := context.Background()
		if telemetryShutdown != nil {
			_ = telemetryShutdown(ctx)
		}
		if pool != nil {
			pool.Close()
		}
	}
}
EOF
  sed -i '' "s|SERVICE_MODULE|$mod|g" "$svc_dir/internal/di/cleanup.go" 2>/dev/null || \
    sed -i "s|SERVICE_MODULE|$mod|g" "$svc_dir/internal/di/cleanup.go"

  write_file "$svc_dir/internal/di/wire_gen.go" <<EOF
// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"$mod/internal/app"
)

// Initialize builds the service application graph.
func Initialize() (*app.App, func(), error) {
	base, err := provideConfig()
	if err != nil {
		return nil, nil, err
	}
	slogLogger := provideLogger(base)
	pool, err := providePool(base)
	if err != nil {
		return nil, nil, err
	}
	queries := provideQueries(pool)
	router := provideRouter(base, slogLogger, queries)
	telemetryShutdown, err := provideTelemetry(base, slogLogger)
	if err != nil {
		return nil, nil, err
	}
	appApp := provideApp(base, slogLogger, router, pool, queries)
	cleanup := newCleanup(appApp, pool, telemetryShutdown)
	return appApp, cleanup, nil
}
EOF

  write_file "$svc_dir/cmd/server/main.go" <<EOF
package main

import (
	"log/slog"
	"os"
	"time"

	"$mod/internal/di"
	"github.com/ai-finops/ai-finops/packages/shared"
)

func main() {
	app, cleanup, err := di.Initialize()
	if err != nil {
		slog.Error("failed to initialize service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cleanup()

	srv := app.HTTPServer()
	if err := shared.ServeWithGracefulShutdown(app.Log, srv, 15*time.Second); err != nil {
		app.Log.Error("server stopped with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
EOF

done

echo "WS5 service scaffolding complete."