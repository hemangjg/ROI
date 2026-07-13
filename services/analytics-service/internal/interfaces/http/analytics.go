package httpapi

import (
	"errors"
	"net/http"

	analyticsapp "github.com/ai-finops/ai-finops/services/analytics-service/internal/application/analytics"
	"github.com/ai-finops/ai-finops/packages/shared"
	"github.com/go-chi/chi/v5"
)

// AnalyticsHandler serves spend analytics REST endpoints.
type AnalyticsHandler struct {
	svc *analyticsapp.Service
}

// NewAnalyticsHandler constructs analytics HTTP handlers.
func NewAnalyticsHandler(svc *analyticsapp.Service) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

// GetSpendSummary handles GET /v1/orgs/{orgId}/spend/summary.
func (h *AnalyticsHandler) GetSpendSummary(w http.ResponseWriter, r *http.Request) {
	params, err := queryParamsFromRequest(r)
	if err != nil {
		writeAnalyticsError(w, err)
		return
	}

	result, err := h.svc.GetSpendSummary(r.Context(), params)
	if err != nil {
		writeAnalyticsError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GetSpendTimeSeries handles GET /v1/orgs/{orgId}/spend/timeseries.
func (h *AnalyticsHandler) GetSpendTimeSeries(w http.ResponseWriter, r *http.Request) {
	params, err := queryParamsFromRequest(r)
	if err != nil {
		writeAnalyticsError(w, err)
		return
	}

	result, err := h.svc.GetSpendTimeSeries(r.Context(), params, r.URL.Query().Get("granularity"))
	if err != nil {
		writeAnalyticsError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GetSpendByProvider handles GET /v1/orgs/{orgId}/spend/by-provider.
func (h *AnalyticsHandler) GetSpendByProvider(w http.ResponseWriter, r *http.Request) {
	params, err := queryParamsFromRequest(r)
	if err != nil {
		writeAnalyticsError(w, err)
		return
	}

	result, err := h.svc.GetSpendByProvider(r.Context(), params)
	if err != nil {
		writeAnalyticsError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GetSpendByTeam handles GET /v1/orgs/{orgId}/spend/by-team.
func (h *AnalyticsHandler) GetSpendByTeam(w http.ResponseWriter, r *http.Request) {
	params, err := queryParamsFromRequest(r)
	if err != nil {
		writeAnalyticsError(w, err)
		return
	}

	result, err := h.svc.GetSpendByTeam(r.Context(), params)
	if err != nil {
		writeAnalyticsError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func queryParamsFromRequest(r *http.Request) (analyticsapp.QueryParams, error) {
	q := r.URL.Query()
	return analyticsapp.ParseQueryParams(
		chi.URLParam(r, "orgId"),
		q.Get("from"),
		q.Get("to"),
		q.Get("team_id"),
	)
}

func writeAnalyticsError(w http.ResponseWriter, err error) {
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		writeJSON(w, appErr.Status, map[string]string{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"code":    "internal_error",
		"message": "internal server error",
	})
}