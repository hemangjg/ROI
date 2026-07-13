package auth_test

import (
	"context"
	"os"
	"testing"
	"time"

	pkgauth "github.com/ai-finops/ai-finops/packages/auth"
	"github.com/ai-finops/ai-finops/packages/logger"
	authapp "github.com/ai-finops/ai-finops/services/identity-service/internal/application/auth"
	jwtkeys "github.com/ai-finops/ai-finops/services/identity-service/internal/infrastructure/jwt"
	"github.com/ai-finops/ai-finops/services/identity-service/internal/infrastructure/postgres"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestAuthServiceIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	privatePath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	publicPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if privatePath == "" {
		privatePath = "../../secrets/jwt/private.pem"
	}
	if publicPath == "" {
		publicPath = "../../secrets/jwt/public.pem"
	}

	keys, err := jwtkeys.Load(privatePath, publicPath)
	if err != nil {
		t.Skipf("jwt keys unavailable: %v", err)
	}

	ctx := context.Background()
	pool, queries, err := postgres.Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer pool.Close()

	svc := authapp.NewService(
		pool,
		queries,
		keys.Private,
		keys.Public,
		15*time.Minute,
		168*time.Hour,
		logger.New("identity-service-test", "error"),
	)

	email := "ws3-" + uuid.NewString() + "@example.com"
	password := "demo-password-change-me"

	tokens, err := svc.Register(ctx, email, password, "WS3 User", "WS3 Org")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" || tokens.OrgID == "" {
		t.Fatal("expected token pair from register")
	}

	valid, userID, orgID, roles, err := svc.ValidateToken(ctx, tokens.AccessToken)
	if err != nil || !valid || userID == "" || orgID != tokens.OrgID || len(roles) == 0 {
		t.Fatalf("validate token: valid=%v user=%q org=%q roles=%v err=%v", valid, userID, orgID, roles, err)
	}

	loginTokens, err := svc.Login(ctx, email, password)
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if loginTokens.AccessToken == "" {
		t.Fatal("expected access token from login")
	}

	plaintext, hash, err := pkgauth.GenerateAPIKey()
	if err != nil {
		t.Fatalf("generate api key: %v", err)
	}
	prefix, err := pkgauth.APIKeyPrefix(plaintext)
	if err != nil {
		t.Fatalf("api key prefix: %v", err)
	}

	_, err = queries.CreateApiKey(ctx, sqlcgen.CreateApiKeyParams{
		OrgID:     mustPgUUID(t, tokens.OrgID),
		Name:      "ws3-test",
		KeyPrefix: prefix,
		KeyHash:   hash,
		Scopes:    []string{"ingest"},
		CreatedBy: mustPgUUID(t, userID),
	})
	if err != nil {
		t.Fatalf("create api key: %v", err)
	}

	apiValid, apiOrgID, keyID, err := svc.ValidateApiKey(ctx, plaintext)
	if err != nil || !apiValid || apiOrgID != tokens.OrgID || keyID == "" {
		t.Fatalf("validate api key: valid=%v org=%q key=%q err=%v", apiValid, apiOrgID, keyID, err)
	}
}

func mustPgUUID(t *testing.T, raw string) pgtype.UUID {
	t.Helper()

	parsed, err := uuid.Parse(raw)
	if err != nil {
		t.Fatalf("parse uuid %q: %v", raw, err)
	}

	return pgtype.UUID{Bytes: parsed, Valid: true}
}