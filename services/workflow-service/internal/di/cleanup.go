package di

import (
	"context"

	"github.com/ai-finops/ai-finops/services/workflow-service/internal/app"
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
		if app != nil {
			if app.KafkaConsumer != nil {
				_ = app.KafkaConsumer.Close()
			}
			if app.PricingClient != nil {
				_ = app.PricingClient.Close()
			}
			if app.ClickHouse != nil {
				_ = app.ClickHouse.Close()
			}
		}
		if pool != nil {
			pool.Close()
		}
	}
}