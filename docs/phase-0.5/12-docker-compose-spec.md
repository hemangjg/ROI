# 12 — Docker Compose Specification

## Profiles

| Profile | Purpose | Services |
|---------|---------|----------|
| `core` | Data + messaging only | postgres, clickhouse, redis, kafka, temporal, openfga |
| `services` | Full backend | core + 6 Go services |
| `web` | With dashboard | services + web |
| `obs` | Observability | otel-collector, prometheus, grafana, jaeger, loki |
| `full` | Everything | web + obs |

### Commands

```bash
# Daily development
docker compose --profile services up

# With dashboard
docker compose --profile web up

# Full stack including observability
docker compose --profile full up
```

---

## Service Definitions

### Infrastructure (core)

| Service | Image | Port | Health check |
|---------|-------|------|-------------|
| postgres | postgres:16-alpine | 5432 | `pg_isready` |
| clickhouse | clickhouse/clickhouse-server:24 | 8123, 9000 | HTTP ping |
| redis | redis:7-alpine | 6379 | `redis-cli ping` |
| kafka | bitnami/kafka:3.7 (KRaft) | 9092 | broker API |
| temporal | temporalio/auto-setup:1.25 | 7233 | gRPC health |
| openfga | openfga/openfga:latest | 8080, 8081 | HTTP /healthz |

### Application (services)

| Service | Build context | Port | Depends on |
|---------|--------------|------|------------|
| identity-service | services/identity-service | 8080, 9090 | postgres, redis, openfga |
| ingestion-service | services/ingestion-service | 8080 | postgres, redis, kafka |
| management-service | services/management-service | 8080 | postgres, openfga |
| analytics-service | services/analytics-service | 8080 | clickhouse, redis, openfga |
| pricing-service | services/pricing-service | 9090 | postgres |
| workflow-service | services/workflow-service | — | kafka, temporal, postgres, clickhouse, pricing-service |

### Frontend (web)

| Service | Build context | Port |
|---------|--------------|------|
| web | apps/web | 3000 |

Env: `NEXT_PUBLIC_API_URL=http://localhost:8080` (or Envoy proxy on 8888 for local routing).

---

## Local Envoy Gateway (optional)

For local ext_authz testing, add `envoy-gateway` service routing to backend services on port 8888.

---

## Environment Variables

Shared via `.env` (gitignored):

```bash
DATABASE_URL=postgres://finops:finops@postgres:5432/finops?sslmode=disable
CLICKHOUSE_URL=clickhouse://clickhouse:9000/finops
REDIS_URL=redis://redis:6379
KAFKA_BROKERS=kafka:9092
TEMPORAL_HOST=temporal:7233
OPENFGA_URL=http://openfga:8080
JWT_PRIVATE_KEY_PATH=/secrets/jwt.pem
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
```

---

## Initialization

### `scripts/init-local.sh`

1. `atlas migrate apply --env local`
2. Apply ClickHouse migrations
3. `buf generate`
4. Seed: demo org, admin user, OpenAI + Anthropic pricing
5. Generate test API key
6. Insert 1000 synthetic events via SDK script

### Seed credentials (local only)

```
Email: admin@demo.ai-finops.local
Password: demo-password-change-me
Org slug: demo
API key: (generated, printed once)
```

---

## Volume Mounts

| Service | Volume | Purpose |
|---------|--------|---------|
| postgres | `pgdata` | Persistent DB |
| clickhouse | `chdata` | Persistent analytics |
| kafka | `kafkadata` | Log retention |

Application services: no persistent volumes (stateless).

---

## Resource Limits (local)

| Service | CPU | Memory |
|---------|-----|--------|
| clickhouse | 2 | 2GB |
| kafka | 1 | 1GB |
| Go services | 0.5 | 256MB each |
| web | 1 | 512MB |

Total recommended: 8GB RAM, 4 CPU cores.