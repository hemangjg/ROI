package httpapi

import (
	"log/slog"
	"net/http"

	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/ai-finops/ai-finops/packages/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter wires health and auth routes.
func NewRouter(serviceName string, log *slog.Logger, queries *sqlcgen.Queries, auth *AuthHandler) chi.Router {
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

	if auth != nil {
		r.Route("/v1/auth", func(r chi.Router) {
			r.Post("/register", auth.Register)
			r.Post("/login", auth.Login)
			r.Post("/refresh", auth.Refresh)
			r.Post("/logout", auth.Logout)
		})
	}

	return r
}