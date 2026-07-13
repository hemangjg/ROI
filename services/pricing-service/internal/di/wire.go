//go:build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/app"
)

func Initialize() (*app.App, func(), error) {
	wire.Build(
		provideConfig,
		provideLogger,
		providePool,
		provideQueries,
		providePluginRegistry,
		providePricingService,
		provideGRPCServer,
		provideRouter,
		provideTelemetry,
		provideApp,
		newCleanup,
	)
	return nil, nil, nil
}