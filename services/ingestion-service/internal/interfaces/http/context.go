package httpapi

import (
	"context"

	ingestapp "github.com/ai-finops/ai-finops/services/ingestion-service/internal/application/ingestion"
)

type apiKeyContextKey struct{}

// WithAPIKeyIdentity stores the authenticated API key on the request context.
func WithAPIKeyIdentity(ctx context.Context, identity ingestapp.APIKeyIdentity) context.Context {
	return context.WithValue(ctx, apiKeyContextKey{}, identity)
}

// APIKeyIdentityFromContext returns the authenticated API key identity.
func APIKeyIdentityFromContext(ctx context.Context) (ingestapp.APIKeyIdentity, bool) {
	identity, ok := ctx.Value(apiKeyContextKey{}).(ingestapp.APIKeyIdentity)
	return identity, ok
}