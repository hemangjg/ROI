package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	ingestapp "github.com/ai-finops/ai-finops/services/ingestion-service/internal/application/ingestion"
	"github.com/ai-finops/ai-finops/packages/shared"
	"github.com/go-chi/chi/v5"
)

// IngestionHandler serves REST ingestion endpoints.
type IngestionHandler struct {
	svc *ingestapp.Service
}

// NewIngestionHandler constructs ingestion HTTP handlers.
func NewIngestionHandler(svc *ingestapp.Service) *IngestionHandler {
	return &IngestionHandler{svc: svc}
}

type batchIngestRequest struct {
	Events []ingestapp.UsageEvent `json:"events"`
}

// IngestEvent handles POST /v1/events.
func (h *IngestionHandler) IngestEvent(w http.ResponseWriter, r *http.Request) {
	identity, ok := requireIdentity(w, r)
	if !ok {
		return
	}

	var event ingestapp.UsageEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeIngestError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	result, err := h.svc.IngestEvent(r.Context(), identity, event)
	if err != nil {
		writeIngestError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, result)
}

// IngestBatch handles POST /v1/events/batch.
func (h *IngestionHandler) IngestBatch(w http.ResponseWriter, r *http.Request) {
	identity, ok := requireIdentity(w, r)
	if !ok {
		return
	}

	var req batchIngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeIngestError(w, shared.BadRequest("invalid_request", "invalid json body"))
		return
	}

	result, err := h.svc.IngestBatch(r.Context(), identity, req.Events)
	if err != nil {
		writeIngestError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, result)
}

// GetEventStatus handles GET /v1/events/{eventId}.
func (h *IngestionHandler) GetEventStatus(w http.ResponseWriter, r *http.Request) {
	identity, ok := requireIdentity(w, r)
	if !ok {
		return
	}

	status, err := h.svc.GetEventStatus(r.Context(), identity, chi.URLParam(r, "eventId"))
	if err != nil {
		writeIngestError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func writeIngestError(w http.ResponseWriter, err error) {
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