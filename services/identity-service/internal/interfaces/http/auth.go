package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	authapp "github.com/ai-finops/ai-finops/services/identity-service/internal/application/auth"
	"github.com/ai-finops/ai-finops/packages/shared"
)

// AuthHandler serves REST auth endpoints.
type AuthHandler struct {
	svc *authapp.Service
}

// NewAuthHandler constructs auth HTTP handlers.
func NewAuthHandler(svc *authapp.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	OrgName  string `json:"org_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Register handles POST /v1/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	tokens, err := h.svc.Register(r.Context(), req.Email, req.Password, req.Name, req.OrgName)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, tokens)
}

// Login handles POST /v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	tokens, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

// Refresh handles POST /v1/auth/refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	tokens, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

// Logout handles POST /v1/auth/logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		writeAuthError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeAuthError(w http.ResponseWriter, err error) {
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