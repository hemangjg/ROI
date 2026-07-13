package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	svcconfig "github.com/ai-finops/ai-finops/services/analytics-service/internal/config"
	chinfra "github.com/ai-finops/ai-finops/services/analytics-service/internal/infrastructure/clickhouse"
	grpcapi "github.com/ai-finops/ai-finops/services/analytics-service/internal/interfaces/grpc"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App holds runtime dependencies for the analytics service.
type App struct {
	Config     svcconfig.Config
	Log        *slog.Logger
	Router     chi.Router
	Pool       *pgxpool.Pool
	Queries    *sqlcgen.Queries
	GRPCServer *grpcapi.Server
	ClickHouse *chinfra.Client
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