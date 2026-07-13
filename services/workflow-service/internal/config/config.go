package config

import "github.com/ai-finops/ai-finops/packages/config"

// Config extends shared service configuration with workflow-service settings.
type Config struct {
	config.Base `koanf:",squash"`
	GRPCPort            int    `koanf:"grpc_port" validate:"required,min=1,max=65535"`
	KafkaBrokers        string `koanf:"kafka_brokers" validate:"required"`
	KafkaTopic          string `koanf:"kafka_topic" validate:"required"`
	KafkaGroupID        string `koanf:"kafka_group_id" validate:"required"`
	PricingGRPCAddr     string `koanf:"pricing_grpc_addr" validate:"required"`
	ClickHouseURL       string `koanf:"clickhouse_url" validate:"required"`
	IdempotencyTTLHours int    `koanf:"idempotency_ttl_hours" validate:"required,min=1"`
}