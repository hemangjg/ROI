package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpapi "github.com/ai-finops/ai-finops/services/identity-service/internal/interfaces/http"
)

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	httpapi.Healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}