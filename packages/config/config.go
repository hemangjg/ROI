package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var validate = validator.New()

// Base holds configuration shared by every service.
type Base struct {
	ServiceName string `koanf:"service_name" validate:"required"`
	Environment string `koanf:"environment" validate:"required,oneof=local development staging production"`
	LogLevel    string `koanf:"log_level" validate:"required,oneof=debug info warn error"`
	HTTPPort    int    `koanf:"http_port" validate:"required,min=1,max=65535"`
	DatabaseURL string `koanf:"database_url"`
	RedisURL    string `koanf:"redis_url"`
	OTLPEndpoint string `koanf:"otlp_endpoint"`
}

// Load reads YAML defaults + optional local overrides + environment variables.
func Load[T any](serviceName string, paths ...string) (T, error) {
	var zero T

	k := koanf.New(".")
	if err := k.Load(file.Provider(defaultsPath(paths)), yaml.Parser()); err != nil {
		return zero, fmt.Errorf("load defaults: %w", err)
	}

	local := localPath(paths)
	if local != "" {
		_ = k.Load(file.Provider(local), yaml.Parser())
	}

	envPrefix := strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_"))
	if err := k.Load(env.Provider(envPrefix+"_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, envPrefix+"_"))
	}), nil); err != nil {
		return zero, fmt.Errorf("load env: %w", err)
	}

	if err := k.Set("service_name", serviceName); err != nil {
		return zero, fmt.Errorf("set service_name: %w", err)
	}

	var cfg T
	if err := k.Unmarshal("", &cfg); err != nil {
		return zero, fmt.Errorf("unmarshal: %w", err)
	}

	if err := validate.Struct(cfg); err != nil {
		return zero, fmt.Errorf("validate: %w", err)
	}

	return cfg, nil
}

func defaultsPath(paths []string) string {
	if len(paths) > 0 && paths[0] != "" {
		return paths[0]
	}
	return "config/defaults.yaml"
}

func localPath(paths []string) string {
	if len(paths) > 1 {
		return paths[1]
	}
	return "config/config.local.yaml"
}