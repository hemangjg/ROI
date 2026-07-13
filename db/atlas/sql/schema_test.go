package sqlcgen_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestSchemaTablesExist(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set — skipping integration test")
	}

	expected := []string{
		"organizations",
		"users",
		"org_memberships",
		"teams",
		"team_members",
		"api_keys",
		"refresh_tokens",
		"providers",
		"models",
		"model_pricing",
		"outbox_events",
		"idempotency_keys",
		"audit_logs",
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	for _, table := range expected {
		var exists bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.tables
				WHERE table_schema = 'public' AND table_name = $1
			)
		`, table).Scan(&exists)
		if err != nil {
			t.Fatalf("query %s: %v", table, err)
		}
		if !exists {
			t.Fatalf("table %s does not exist — run atlas migrate apply", table)
		}
	}
}