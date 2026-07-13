//go:build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/ai-finops/ai-finops/services/workflow-service/internal/app"
)

func Initialize() (*app.App, func(), error) {
	wire.Build(
		provideConfig,
		provideLogger,
		providePool,
		provideQueries,
		providePricingClient,
		provideClickHouseClient,
		provideWorkflowService,
		provideKafkaConsumer,
		provideGRPCServer,
		provideRouter,
		provideTelemetry,
		provideApp,
		newCleanup,
	)
	return nil, nil, nil
}