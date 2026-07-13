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