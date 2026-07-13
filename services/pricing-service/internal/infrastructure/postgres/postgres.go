package postgres

import (
	"context"
	"fmt"

	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Open creates a PostgreSQL connection pool and SQLC queries handle.
func Open(ctx context.Context, databaseURL string) (*pgxpool.Pool, *sqlcgen.Queries, error) {
	if databaseURL == "" {
		return nil, nil, fmt.Errorf("database_url is required")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("connect postgres: %w", err)
	}

	return pool, sqlcgen.New(pool), nil
}