package ingestion

import (
	"context"
	"errors"
	"strings"

	pkgauth "github.com/ai-finops/ai-finops/packages/auth"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// AuthenticateAPIKey validates a bearer API key and returns org/key identity.
func AuthenticateAPIKey(ctx context.Context, queries *sqlcgen.Queries, rawKey string) (APIKeyIdentity, error) {
	rawKey = stripBearer(strings.TrimSpace(rawKey))
	if rawKey == "" {
		return APIKeyIdentity{}, shared.Unauthorized("missing_api_key", "api key is required")
	}

	prefix, err := pkgauth.APIKeyPrefix(rawKey)
	if err != nil {
		return APIKeyIdentity{}, shared.Unauthorized("invalid_api_key", "api key is invalid")
	}

	row, err := queries.GetApiKeyByPrefix(ctx, prefix)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return APIKeyIdentity{}, shared.Unauthorized("invalid_api_key", "api key is invalid")
		}
		return APIKeyIdentity{}, shared.Internal("lookup api key")
	}

	ok, err := pkgauth.VerifyAPIKey(rawKey, row.KeyHash)
	if err != nil || !ok {
		return APIKeyIdentity{}, shared.Unauthorized("invalid_api_key", "api key is invalid")
	}

	if !hasIngestScope(row.Scopes) {
		return APIKeyIdentity{}, shared.Forbidden("missing_scope", "api key is not authorized for ingestion")
	}

	return APIKeyIdentity{
		OrgID: pgUUIDString(row.OrgID),
		KeyID: pgUUIDString(row.ID),
	}, nil
}

func hasIngestScope(scopes []string) bool {
	for _, scope := range scopes {
		if scope == "ingest" {
			return true
		}
	}
	return false
}

func stripBearer(token string) string {
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return strings.TrimSpace(token[7:])
	}
	return token
}

func pgUUIDString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	u, err := uuid.FromBytes(id.Bytes[:])
	if err != nil {
		return ""
	}
	return u.String()
}

func parseUUID(raw string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return pgtype.UUID{}, shared.BadRequest("invalid_uuid", "invalid id format")
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}