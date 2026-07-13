package sqlcgen_test

import (
	"context"
	"os"
	"testing"

	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPingIntegration(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set — skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	q := sqlcgen.New(pool)
	ok, err := q.Ping(ctx)
	if err != nil {
		t.Fatalf("ping: %v", err)
	}
	if ok != 1 {
		t.Fatalf("ok = %d, want 1", ok)
	}
}