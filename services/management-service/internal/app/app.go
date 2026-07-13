package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	svcconfig "github.com/ai-finops/ai-finops/services/management-service/internal/config"
	grpcapi "github.com/ai-finops/ai-finops/services/management-service/internal/interfaces/grpc"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App holds runtime dependencies for the service shell.
type App struct {
	Config     svcconfig.Config
	Log        *slog.Logger
	Router     chi.Router
	Pool       *pgxpool.Pool
	Queries    *sqlcgen.Queries
	GRPCServer *grpcapi.Server
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