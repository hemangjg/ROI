# AI FinOps — agent / contributor guide

## What this is

Enterprise AI FinOps platform monorepo: track LLM spend across providers.
Positioning: “Datadog for AI spend.”

## Working directory

```
/Users/hemangjg/Hemangs stuff/Projects/ai-finops
```

Quote all paths (spaces in `Hemangs stuff`).

## Stack (locked)

- Nx + pnpm monorepo
- Go services (Chi, SQLC, Wire)
- Next.js dashboard
- PostgreSQL + ClickHouse + Kafka + Redis + Temporal + OpenFGA + Envoy
- Protos via Buf; REST OpenAPI at `api/openapi/v1.yaml`

## Quality bar (before features)

A change is not done unless:

1. `bash scripts/quality-gate.sh` passes (or at least `bash scripts/ci-local.sh`)
2. For pipeline/Envoy/Kafka changes: `WITH_E2E=1 bash scripts/quality-gate.sh`
3. No secrets committed (`secrets/jwt/`, `*.pem`, `.env`)
4. Docs match reality (`README.md`, `docs/status.md`)

## Commands

```bash
# Offline CI (lint, test, build)
bash scripts/quality-gate.sh

# Live end-to-end ingest → spend
docker compose --profile services up -d
bash scripts/init-local.sh
bash scripts/e2e-smoke.sh

# Or combined:
WITH_E2E=1 bash scripts/quality-gate.sh
```

## Service map

| Service | Host port | Role |
|---------|-----------|------|
| identity | 8080 / gRPC 9090 | auth + ext_authz |
| ingestion | 8081 | usage events + outbox write |
| management | 8082 | orgs, teams, API keys |
| analytics | 8083 | spend queries (ClickHouse) |
| pricing | 8084 / gRPC 9092 | catalog + CalculateCost |
| workflow | 8085 | Kafka consumer → price → CH |
| outbox-relay | — | outbox → Kafka `usage.events.raw` |
| Envoy | 8888 | API gateway |

## Pipeline invariants

1. Ingest writes Postgres outbox, returns 202
2. Outbox-relay publishes protobuf `UsageEventCreated` to `usage.events.raw`
3. Workflow ensures topic exists on boot, then prices + inserts ClickHouse
4. Analytics reads `finops.daily_spend_rollups`
5. Envoy routes `/v1/orgs/{id}/spend*` to analytics (regex is **full-path** match)

## Known footguns

- Path has spaces — always quote
- Kafka consumer can get 0 partitions if topic missing at first join — `EnsureTopic` + `init-local.sh` mitigate
- Nullable UUID → ClickHouse must use untyped nil, not `*uuid.UUID(nil)`
- Envoy `safe_regex` matches the **entire** path
- `packages/proto/gen/` is **gitignored** — CI and local lint must run `pnpm proto:generate` before golangci-lint / go test

## Where to add Phase 2 features

Only after quality gate is green:

- Budgets / alerts: new tables in `db/atlas`, management or new service, dashboard pages
- Extra providers: plugins under `services/pricing-service/internal/plugin/`

Do not expand scope without updating OpenAPI + protos + tests.
