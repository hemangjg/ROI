package di

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ai-finops/ai-finops/services/identity-service/internal/app"
)

func newCleanup(
	app *app.App,
	pool *pgxpool.Pool,
	telemetryShutdown func(context.Context) error,
) func() {
	return func() {
		ctx := context.Background()
		if telemetryShutdown != nil {
			_ = telemetryShutdown(ctx)
		}
		if pool != nil {
			pool.Close()
		}
	}
}