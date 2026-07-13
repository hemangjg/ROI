package config

import "github.com/ai-finops/ai-finops/packages/config"

// Config extends shared service configuration with identity-service settings.
type Config struct {
	config.Base `koanf:",squash"`
	GRPCPort           int    `koanf:"grpc_port" validate:"required,min=1,max=65535"`
	JWTPrivateKeyPath  string `koanf:"jwt_private_key_path" validate:"required"`
	JWTPublicKeyPath   string `koanf:"jwt_public_key_path" validate:"required"`
	AccessTokenTTL     string `koanf:"access_token_ttl" validate:"required"`
	RefreshTokenTTL    string `koanf:"refresh_token_ttl" validate:"required"`
}