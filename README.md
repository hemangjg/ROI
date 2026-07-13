# AI FinOps Platform

Enterprise AI FinOps platform for tracking, governing, and optimizing LLM spending across every provider.

**The Datadog for AI spend** — cross-provider visibility, FinOps governance, and cost optimization in one platform.

---

## Project Location

```
/Users/hemangjg/Hemangs stuff/Projects/ai-finops
```

All project files live at this path on your machine (`~/Hemangs stuff/Projects/ai-finops`).

---

## Current Phase

**Phase 1b — Ingestion MVP** (engineering complete — re-verify with setup scripts + Docker)

| Workstream | Status |
|------------|--------|
| P2-WS1 PostgreSQL + ClickHouse schema | Complete (run `bash scripts/setup-p2-ws1.sh`) |
| P2-WS2 Domain protos | Complete (run `bash scripts/setup-p2-ws2.sh`) |
| P2-WS3 identity-service auth APIs | Complete (run `bash scripts/setup-p2-ws3.sh`) |
| P2-WS4 management-service APIs | Complete (run `bash scripts/setup-p2-ws4.sh`) |
| P2-WS5 pricing-service + catalog | Complete (run `bash scripts/setup-p2-ws5.sh`) |
| P2-WS6 ingestion-service + outbox | Complete (run `bash scripts/setup-p2-ws6.sh`) |
| P2-WS7 outbox relay + Kafka + workflow | Complete (run `bash scripts/setup-p2-ws7.sh`) |
| P2-WS8 analytics-service | Complete (run `bash scripts/setup-p2-ws8.sh`) |
| P2-WS9 Envoy ext_authz + dashboard live APIs | Complete (run `bash scripts/setup-p2-ws9.sh`) |

**Phase 1 — Project Foundation** (complete)

| Workstream | Status |
|------------|--------|
| WS1 Monorepo Bootstrap | Complete |
| WS2 Proto Toolchain | Complete (run `bash scripts/setup-ws2.sh` to verify buf) |
| WS3 DB Foundation | Complete (run `bash scripts/setup-ws3.sh` to verify atlas/sqlc) |
| WS4 Shared Packages | Complete |
| WS5 Go Service Shells | Complete (run `bash scripts/setup-ws5.sh` to verify) |
| WS6 Docker Compose | Complete (run `bash scripts/setup-ws6.sh` to verify) |
| WS7 OpenTelemetry | Complete (run `bash scripts/setup-ws7.sh` to verify) |
| WS8 Envoy Gateway | Complete (run `bash scripts/setup-ws8.sh` to verify) |
| WS9 Frontend Dashboard | Complete (run `bash scripts/setup-ws9.sh` to verify) |
| WS10 TypeScript SDK | Complete (run `bash scripts/setup-ws10.sh` to verify) |
| WS11 Dev UX + CI | Complete (run `bash scripts/setup-ws11.sh` to verify) |
| WS12 Docs + Final Verification | Complete (run `bash scripts/setup-ws12.sh` to verify) |

### Quick start

```bash
pnpm install
go work sync
pnpm graph

# Native dev (infra in Docker, Go services + dashboard on host)
pnpm dev

# Full Docker stack (containerized services)
pnpm dev:docker
# or: make dev

# Verify toolchains (run each script separately — do not paste inline comments)
bash scripts/setup-ws2.sh
bash scripts/setup-ws3.sh

# Test all Go packages (go.work modules — not `go test ./packages/...`)
bash scripts/test-packages.sh

# Build and test all Go service shells
bash scripts/setup-ws5.sh

# Start local infrastructure + services (Docker Compose)
bash scripts/setup-ws6.sh
# or: make dev

# Verify OpenTelemetry (obs stack + traced requests in Jaeger)
bash scripts/setup-ws7.sh
# or: OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317 docker compose --profile full up --build

# Verify Envoy Gateway (local routing + TLS on :8888/:8443)
bash scripts/setup-ws8.sh
```

### API Gateway (Envoy)

| Listener | URL | Purpose |
|----------|-----|---------|
| HTTP | http://localhost:8888/v1 | Local API (matches OpenAPI `servers`) |
| HTTPS API | https://api.ai-finops.local:8443/v1 | TLS API (self-signed cert via `make certs`) |
| HTTPS App | https://app.ai-finops.local:8444 | Dashboard via Envoy (`web` profile) |

Route map: `/v1/events` → ingestion, `/v1/auth` → identity, `/v1/orgs` → management, `/v1/orgs/*/spend` → analytics.

ext_authz validates JWTs via identity-service gRPC (`:9090`) on `/v1/orgs` and `/v1/orgs/*/spend`. `/v1/auth` and `/v1/events` bypass auth.

### Frontend dashboard

```bash
bash scripts/setup-ws9.sh
# or: cd apps/web && pnpm dev
```

Dashboard: http://localhost:3000 (set `NEXT_PUBLIC_USE_MOCK_API=false` for live APIs via Envoy).

### TypeScript SDK

```bash
bash scripts/setup-ws10.sh
```

```typescript
import { FinOpsClient } from "@ai-finops/sdk";

const client = new FinOpsClient({
  apiKey: process.env.AI_FINOPS_API_KEY,
  baseUrl: "http://localhost:8888/v1",
});

await client.ingest({
  provider: "openai",
  model: "gpt-4o",
  input_tokens: 1200,
  output_tokens: 340,
  occurred_at: new Date().toISOString(),
});
```

`client.track(...)` is a fire-and-forget alias for `ingest`. Auth helpers: `login`, `register`, `refresh`.

### CI and quality

```bash
bash scripts/setup-ws11.sh
# or: pnpm ci:local

# Skip flaky host Next.js build locally (CI still builds web via Docker)
SKIP_WEB_BUILD=1 bash scripts/ci-local.sh
```

GitHub Actions workflow: `.github/workflows/ci.yml` (lint, proto, atlas, sqlc, test, build, web image).

Root `pnpm build` excludes the workspace meta-target to avoid recursive `nx run-many` loops.

### Phase 1 sign-off

```bash
bash scripts/setup-ws12.sh
# or: make ws12
```

C4 diagrams: [docs/architecture/](./docs/architecture/).

### Phase 1b — schema (P2-WS1)

```bash
bash scripts/setup-p2-ws1.sh
# or: make p2-ws1
```

Applies PostgreSQL business schema (Atlas), ClickHouse analytics tables, and SQLC queries.

### Docker Compose profiles

| Profile | Command | Includes |
|---------|---------|----------|
| `core` | `docker compose --profile core up -d` | Postgres, ClickHouse, Redis, Kafka, Temporal, OpenFGA |
| `services` | `docker compose --profile services up --build` | core + 6 Go services |
| `web` | `docker compose --profile web up --build` | services + dashboard placeholder |
| `obs` | `docker compose --profile obs up -d` | Jaeger, Loki, Prometheus, Grafana, OTel collector |
| `full` | `docker compose --profile full up --build` | everything |

Stop any standalone `finops-pg` container before starting compose (both use port 5432).

### Service ports (local)

| Service | Port |
|---------|------|
| identity-service | 8080 |
| ingestion-service | 8081 |
| management-service | 8082 |
| analytics-service | 8083 |
| pricing-service | 8084 |
| workflow-service | 8085 |

Run a single service:

```bash
cd services/identity-service
go run ./cmd/server
curl localhost:8080/healthz
```

### Documentation

| Path | Contents |
|------|----------|
| [docs/phase-0/](./docs/phase-0/) | 12 product definition chapters |
| [docs/phase-0.5/](./docs/phase-0.5/) | 12 technical architecture chapters |
| [docs/architecture/](./docs/architecture/) | C4 context, container, component diagrams |
| [docs/adr/](./docs/adr/) | 8 Architecture Decision Records |
| [api/openapi/v1.yaml](./api/openapi/v1.yaml) | External REST API contract |
| [contracts/proto/](./contracts/proto/) | Active Buf module (bootstrap only) |
| [docs/reference/phase-0.5-proto/](./docs/reference/phase-0.5-proto/) | Archived Phase 0.5 protos |
| [db/atlas/](./db/atlas/) | PostgreSQL Atlas + SQLC (empty schema) |

---

## Service Topology

```
Clients (SDKs + Dashboard)
        → Cloudflare → Envoy Gateway
        → identity-service + OpenFGA
        → Ingestion | Management | Analytics | Pricing | Workflow
        → Kafka + Temporal
        → PostgreSQL | ClickHouse | Redis | OpenSearch | S3
        → OpenTelemetry → Prometheus → Grafana | Jaeger | Loki
```

---

## Locked Tech Stack

| Layer | Technology |
|-------|-----------|
| Monorepo | Nx + pnpm |
| Backend | Go 1.25+, Chi, Wire, SQLC, Koanf, slog |
| Frontend | Next.js, shadcn/ui, Tailwind, TanStack, ECharts |
| Business DB | PostgreSQL 16 (Atlas migrations) |
| Analytics DB | ClickHouse |
| Cache | Redis |
| Streaming | Kafka |
| Workflows | Temporal |
| AuthZ | OpenFGA |
| Internal RPC | gRPC + Buf + Protocol Buffers |
| External API | REST + OpenAPI 3.1 |
| Edge | Cloudflare + Envoy Gateway |
| Cloud | AWS (EKS, RDS, MSK, ElastiCache, S3) |
| IaC | Terraform + Helm |
| Observability | OpenTelemetry, Prometheus, Grafana, Jaeger, Loki |

---

## Roadmap

| Phase | Focus | Status |
|-------|-------|--------|
| **0** | Product definition | Complete |
| **0.5** | Technical architecture | Complete |
| **1** | Project foundation (monorepo, shells, CI, docs) | Complete |
| **1b** | Ingestion MVP (domain APIs + data model) | Complete (verify with Docker + `scripts/setup-p2-ws*.sh`) |
| **2** | Budgets + governance | Not started |
| **3** | Optimization + ROI | Not started |
| **4** | Enterprise (SSO, SCIM, billing) | Not started |

---

## Getting Started

1. Read [Phase 0 docs](./docs/phase-0/README.md) for product context
2. Read [Phase 0.5 docs](./docs/phase-0.5/README.md) for architecture
3. Review [ADRs](./docs/adr/) for key decisions
4. Phase 1 foundation: `bash scripts/setup-ws12.sh` then `pnpm dev`
5. Phase 1b: `bash scripts/setup-p2-ws1.sh` then continue with protos and business APIs
6. End-to-end smoke (Docker services up): `bash scripts/e2e-smoke.sh`
7. Verification log: [docs/status.md](./docs/status.md)