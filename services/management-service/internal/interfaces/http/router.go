package httpapi

import (
	"log/slog"
	"net/http"

	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/ai-finops/ai-finops/packages/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter wires health and management routes.
func NewRouter(serviceName string, log *slog.Logger, queries *sqlcgen.Queries, mgmt *ManagementHandler) chi.Router {
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

	if mgmt != nil {
		r.Route("/v1/orgs/{orgId}", func(r chi.Router) {
			r.Get("/", mgmt.GetOrg)
			r.Patch("/", mgmt.PatchOrg)

			r.Route("/teams", func(r chi.Router) {
				r.Post("/", mgmt.CreateTeam)
				r.Get("/", mgmt.ListTeams)
				r.Delete("/{teamId}", mgmt.DeleteTeam)
			})

			r.Route("/users", func(r chi.Router) {
				r.Post("/invite", mgmt.InviteUser)
				r.Get("/", mgmt.ListUsers)
			})

			r.Route("/api-keys", func(r chi.Router) {
				r.Post("/", mgmt.CreateApiKey)
				r.Get("/", mgmt.ListApiKeys)
				r.Delete("/{keyId}", mgmt.RevokeApiKey)
			})

			r.Route("/budgets", func(r chi.Router) {
				r.Get("/", mgmt.ListBudgets)
				r.Put("/", mgmt.UpsertBudget)
				r.Post("/evaluate", mgmt.EvaluateBudget)
			})

			r.Route("/budget-alerts", func(r chi.Router) {
				r.Get("/", mgmt.ListBudgetAlerts)
				r.Post("/{alertId}/ack", mgmt.AcknowledgeBudgetAlert)
			})

			r.Get("/audit-logs", mgmt.ListAuditLogs)
		})
	}

	return r
}