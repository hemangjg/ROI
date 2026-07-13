module github.com/ai-finops/ai-finops/services/outbox-relay

go 1.25.0

require (
	github.com/ai-finops/ai-finops/db/atlas/sql v0.0.0
	github.com/ai-finops/ai-finops/packages/config v0.0.0
	github.com/ai-finops/ai-finops/packages/events v0.0.0
	github.com/ai-finops/ai-finops/packages/logger v0.0.0
	github.com/ai-finops/ai-finops/packages/proto v0.0.0
	github.com/jackc/pgx/v5 v5.7.4
	github.com/segmentio/kafka-go v0.4.47
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.25.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/parsers/yaml v0.1.0 // indirect
	github.com/knadh/koanf/providers/env v1.0.0 // indirect
	github.com/knadh/koanf/providers/file v1.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.27 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/sdk v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/grpc v1.82.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/ai-finops/ai-finops/db/atlas/sql => ../../db/atlas/sql
	github.com/ai-finops/ai-finops/packages/config => ../../packages/config
	github.com/ai-finops/ai-finops/packages/events => ../../packages/events
	github.com/ai-finops/ai-finops/packages/logger => ../../packages/logger
	github.com/ai-finops/ai-finops/packages/proto => ../../packages/proto
)
