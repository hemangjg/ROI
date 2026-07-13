package httpapi

import (
	"net/http"

	ingestapp "github.com/ai-finops/ai-finops/services/ingestion-service/internal/application/ingestion"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
)

// APIKeyAuth validates bearer API keys for ingestion routes.
func APIKeyAuth(queries *sqlcgen.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity, err := ingestapp.AuthenticateAPIKey(r.Context(), queries, r.Header.Get("Authorization"))
			if err != nil {
				writeIngestError(w, err)
				return
			}
			ctx := WithAPIKeyIdentity(r.Context(), identity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func requireIdentity(w http.ResponseWriter, r *http.Request) (ingestapp.APIKeyIdentity, bool) {
	identity, ok := APIKeyIdentityFromContext(r.Context())
	if !ok {
		writeIngestError(w, shared.Unauthorized("missing_api_key", "api key is required"))
		return ingestapp.APIKeyIdentity{}, false
	}
	return identity, true
}