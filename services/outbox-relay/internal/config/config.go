package config

import "github.com/ai-finops/ai-finops/packages/config"

// Config holds outbox-relay settings.
type Config struct {
	config.Base `koanf:",squash"`
	KafkaBrokers  string `koanf:"kafka_brokers" validate:"required"`
	KafkaTopic    string `koanf:"kafka_topic" validate:"required"`
	PollInterval  string `koanf:"poll_interval" validate:"required"`
	BatchSize     int    `koanf:"batch_size" validate:"required,min=1,max=500"`
}