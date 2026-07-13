package postgres

import (
	"context"
	"fmt"

	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Open connects to PostgreSQL and returns a pool plus SQLC queries.
func Open(ctx context.Context, databaseURL string) (*pgxpool.Pool, *sqlcgen.Queries, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, sqlcgen.New(pool), nil
}