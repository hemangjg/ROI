package di

import (
	"context"

	"github.com/ai-finops/ai-finops/services/analytics-service/internal/app"
	"github.com/jackc/pgx/v5/pgxpool"
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
		if app.ClickHouse != nil {
			_ = app.ClickHouse.Close()
		}
		if pool != nil {
			pool.Close()
		}
	}
}