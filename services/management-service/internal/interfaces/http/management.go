package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	mgmtapp "github.com/ai-finops/ai-finops/services/management-service/internal/application/management"
	"github.com/ai-finops/ai-finops/packages/shared"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ManagementHandler serves REST management endpoints.
type ManagementHandler struct {
	svc *mgmtapp.Service
}

// NewManagementHandler constructs management HTTP handlers.
func NewManagementHandler(svc *mgmtapp.Service) *ManagementHandler {
	return &ManagementHandler{svc: svc}
}

type patchOrgRequest struct {
	Name     *string `json:"name"`
	Timezone *string `json:"timezone"`
}

type createTeamRequest struct {
	Name string `json:"name"`
}

type inviteUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type createApiKeyRequest struct {
	Name string `json:"name"`
}

type upsertBudgetRequest struct {
	TeamID           string `json:"team_id"`
	AmountUSD        string `json:"amount_usd"`
	SoftThresholdPct *int32 `json:"soft_threshold_pct,omitempty"`
	PeriodStart      string `json:"period_start,omitempty"`
}

type evaluateBudgetRequest struct {
	TeamID      string `json:"team_id"`
	SpendUSD    string `json:"spend_usd"`
	PeriodStart string `json:"period_start,omitempty"`
}

// GetOrg handles GET /v1/orgs/{orgId}.
func (h *ManagementHandler) GetOrg(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	org, err := h.svc.GetOrg(r.Context(), orgID, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, org)
}

// PatchOrg handles PATCH /v1/orgs/{orgId}.
func (h *ManagementHandler) PatchOrg(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	var req patchOrgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeMgmtError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	org, err := h.svc.UpdateOrg(r.Context(), orgID, req.Name, req.Timezone, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, org)
}

// CreateTeam handles POST /v1/orgs/{orgId}/teams.
func (h *ManagementHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	var req createTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeMgmtError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	team, err := h.svc.CreateTeam(r.Context(), orgID, req.Name, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, team)
}

// ListTeams handles GET /v1/orgs/{orgId}/teams.
func (h *ManagementHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	teams, err := h.svc.ListTeams(r.Context(), orgID, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, teams)
}

// DeleteTeam handles DELETE /v1/orgs/{orgId}/teams/{teamId}.
func (h *ManagementHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	teamID := chi.URLParam(r, "teamId")

	if err := h.svc.DeleteTeam(r.Context(), orgID, teamID, actorFromRequest(r)); err != nil {
		writeMgmtError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// InviteUser handles POST /v1/orgs/{orgId}/users/invite.
func (h *ManagementHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	var req inviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeMgmtError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	member, err := h.svc.InviteUser(r.Context(), orgID, req.Email, req.Name, req.Role, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, member)
}

// ListUsers handles GET /v1/orgs/{orgId}/users.
func (h *ManagementHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	users, err := h.svc.ListUsers(r.Context(), orgID, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// CreateApiKey handles POST /v1/orgs/{orgId}/api-keys.
func (h *ManagementHandler) CreateApiKey(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	var req createApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeMgmtError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	created, err := h.svc.CreateApiKey(r.Context(), orgID, req.Name, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListApiKeys handles GET /v1/orgs/{orgId}/api-keys.
func (h *ManagementHandler) ListApiKeys(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")

	keys, err := h.svc.ListApiKeys(r.Context(), orgID, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, keys)
}

// RevokeApiKey handles DELETE /v1/orgs/{orgId}/api-keys/{keyId}.
func (h *ManagementHandler) RevokeApiKey(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	keyID := chi.URLParam(r, "keyId")

	if err := h.svc.RevokeApiKey(r.Context(), orgID, keyID, actorFromRequest(r)); err != nil {
		writeMgmtError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListAuditLogs handles GET /v1/orgs/{orgId}/audit-logs.
func (h *ManagementHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	cursor := r.URL.Query().Get("cursor")

	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			writeMgmtError(w, shared.BadRequest("invalid_limit", "limit must be an integer"))
			return
		}
		limit = parsed
	}

	page, err := h.svc.ListAuditLogs(r.Context(), orgID, cursor, limit, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

// UpsertBudget handles PUT /v1/orgs/{orgId}/budgets.
func (h *ManagementHandler) UpsertBudget(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	var req upsertBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeMgmtError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}
	var threshold int32 = 80
	if req.SoftThresholdPct != nil {
		threshold = *req.SoftThresholdPct
	}
	budget, err := h.svc.UpsertBudget(r.Context(), orgID, req.TeamID, req.AmountUSD, threshold, req.PeriodStart, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, budget)
}

// ListBudgets handles GET /v1/orgs/{orgId}/budgets.
func (h *ManagementHandler) ListBudgets(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	period := r.URL.Query().Get("period")
	budgets, err := h.svc.ListBudgets(r.Context(), orgID, period, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, budgets)
}

// EvaluateBudget handles POST /v1/orgs/{orgId}/budgets/evaluate.
func (h *ManagementHandler) EvaluateBudget(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	var req evaluateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeMgmtError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}
	status, err := h.svc.EvaluateBudget(r.Context(), orgID, req.TeamID, req.PeriodStart, req.SpendUSD, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

// ListBudgetAlerts handles GET /v1/orgs/{orgId}/budget-alerts.
func (h *ManagementHandler) ListBudgetAlerts(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	status := r.URL.Query().Get("status")
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			writeMgmtError(w, shared.BadRequest("invalid_limit", "limit must be an integer"))
			return
		}
		limit = parsed
	}
	alerts, err := h.svc.ListBudgetAlerts(r.Context(), orgID, status, limit, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, alerts)
}

// AcknowledgeBudgetAlert handles POST /v1/orgs/{orgId}/budget-alerts/{alertId}/ack.
func (h *ManagementHandler) AcknowledgeBudgetAlert(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	alertID := chi.URLParam(r, "alertId")
	alert, err := h.svc.AcknowledgeBudgetAlert(r.Context(), orgID, alertID, actorFromRequest(r))
	if err != nil {
		writeMgmtError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, alert)
}

func actorFromRequest(r *http.Request) mgmtapp.Actor {
	actor := mgmtapp.Actor{}
	if raw := r.Header.Get("X-User-Id"); raw != "" {
		if parsed, err := uuid.Parse(raw); err == nil {
			actor.UserID = pgtype.UUID{Bytes: parsed, Valid: true}
		}
	}
	return actor
}

func writeMgmtError(w http.ResponseWriter, err error) {
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