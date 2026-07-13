package ingestion_test

import (
	"context"
	"os"
	"testing"
	"time"

	pkgauth "github.com/ai-finops/ai-finops/packages/auth"
	"github.com/ai-finops/ai-finops/packages/logger"
	ingestapp "github.com/ai-finops/ai-finops/services/ingestion-service/internal/application/ingestion"
	"github.com/ai-finops/ai-finops/services/ingestion-service/internal/infrastructure/postgres"
	redinfra "github.com/ai-finops/ai-finops/services/ingestion-service/internal/infrastructure/redis"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestIngestionServiceIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, queries, err := postgres.Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer pool.Close()

	limiter, err := redinfra.NewRateLimiter("", 1000, logger.New("ingestion-test", "error"))
	if err != nil {
		t.Fatalf("rate limiter: %v", err)
	}
	defer func() { _ = limiter.Close() }()

	svc := ingestapp.NewService(pool, queries, limiter, logger.New("ingestion-service-test", "error"))

	org, err := queries.CreateOrganization(ctx, sqlcgen.CreateOrganizationParams{
		Name:    "WS6 Org",
		Slug:    "ws6-" + uuid.NewString()[:8],
		Column3: "UTC",
	})
	if err != nil {
		t.Fatalf("create org: %v", err)
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
		OrgID:     org.ID,
		Name:      "ws6-ingest",
		KeyPrefix: prefix,
		KeyHash:   hash,
		Scopes:    []string{"ingest"},
		CreatedBy: pgtype.UUID{},
	})
	if err != nil {
		t.Fatalf("create api key: %v", err)
	}

	identity, err := ingestapp.AuthenticateAPIKey(ctx, queries, "Bearer "+plaintext)
	if err != nil {
		t.Fatalf("authenticate api key: %v", err)
	}

	event := ingestapp.UsageEvent{
		IdempotencyKey: "ws6-" + uuid.NewString(),
		Provider:       "openai",
		Model:          "gpt-4o",
		InputTokens:    1000,
		OutputTokens:   500,
		OccurredAt:     time.Now().UTC(),
	}

	result, err := svc.IngestEvent(ctx, identity, event)
	if err != nil {
		t.Fatalf("ingest event: %v", err)
	}
	if result.EventID == "" || result.Status != "queued" {
		t.Fatalf("unexpected ingest result: %+v", result)
	}

	status, err := svc.GetEventStatus(ctx, identity, result.EventID)
	if err != nil {
		t.Fatalf("get event status: %v", err)
	}
	if status.EventID != result.EventID || status.Status != "queued" {
		t.Fatalf("unexpected status: %+v", status)
	}

	rows, err := queries.ListUnpublishedOutboxEvents(ctx, 10)
	if err != nil {
		t.Fatalf("list outbox: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected unpublished outbox events")
	}
}