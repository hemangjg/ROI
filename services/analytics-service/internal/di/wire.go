//go:build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/ai-finops/ai-finops/services/analytics-service/internal/app"
)

func Initialize() (*app.App, func(), error) {
	wire.Build(
		provideConfig,
		provideLogger,
		providePool,
		provideQueries,
		provideClickHouseClient,
		provideAnalyticsService,
		provideAnalyticsHandler,
		provideGRPCServer,
		provideRouter,
		provideTelemetry,
		provideApp,
		newCleanup,
	)
	return nil, nil, nil
}