//go:build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/ai-finops/ai-finops/services/identity-service/internal/app"
)

func Initialize() (*app.App, func(), error) {
	wire.Build(
		provideConfig,
		provideLogger,
		providePool,
		provideQueries,
		provideJWTKeys,
		provideAuthService,
		provideAuthHandler,
		provideGRPCServer,
		provideRouter,
		provideTelemetry,
		provideApp,
		newCleanup,
	)
	return nil, nil, nil
}