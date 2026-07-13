package config

import "github.com/ai-finops/ai-finops/packages/config"

// Config extends shared service configuration with management-service settings.
type Config struct {
	config.Base `koanf:",squash"`
	GRPCPort int `koanf:"grpc_port" validate:"required,min=1,max=65535"`
}