# Phase 0.5 — Technical Architecture

Architecture documentation for the AI FinOps platform. Defines **how** to build before Phase 1 implementation.

## Locked Service Topology

```
Clients (SDKs + Dashboard)
        → Cloudflare → Envoy Gateway
        → identity-service + OpenFGA
        → Ingestion | Management | Analytics | Pricing | Workflow Services
        → Kafka + Temporal
        → PostgreSQL | ClickHouse | Redis | OpenSearch | S3
        → OpenTelemetry → Prometheus → Grafana | Jaeger | Loki
```

## Reading Order

| # | Document | Contents |
|---|----------|----------|
| 01 | [System Architecture](./01-system-architecture.md) | C4 diagrams, service catalog, Envoy routing |
| 02 | [Data Model](./02-data-model.md) | PostgreSQL (Atlas), ClickHouse, Redis keys |
| 03 | [API Contracts](./03-api-contracts.md) | REST + gRPC (Buf) specifications |
| 04 | [Ingestion Pipeline](./04-ingestion-pipeline.md) | Outbox, Kafka, Temporal workflows |
| 05 | [Security and Tenancy](./05-security-and-tenancy.md) | Auth Service, OpenFGA, encryption |
| 06 | [Frontend Architecture](./06-frontend-architecture.md) | Next.js, shadcn, ECharts, TanStack |
| 07 | [Monorepo Structure](./07-monorepo-structure.md) | Nx + pnpm layout |
| 08 | [AI Provider Plugins](./08-ai-provider-plugins.md) | Plugin interface, provider matrix |
| 09 | [Observability](./09-observability.md) | OTel, Prometheus, SLIs |
| 10 | [Infrastructure](./10-infrastructure.md) | AWS, Terraform, Envoy Gateway, deployment path |
| 11 | [CI/CD and Quality](./11-cicd-and-quality.md) | GitHub Actions, testing, linting |
| 12 | [Docker Compose](./12-docker-compose-spec.md) | Local dev profiles |

## Related Artifacts

| Path | Contents |
|------|----------|
| [../adr/](../adr/) | 8 Architecture Decision Records |
| [../../api/openapi/v1.yaml](../../api/openapi/v1.yaml) | External REST contract |
| [../../api/proto/](../../api/proto/) | Buf module (internal gRPC) |
| [../../db/atlas/](../../db/atlas/) | PostgreSQL declarative schema |

## Exit Criteria

- [x] Service topology documented (6 services + Auth)
- [x] Dual-store data model (PostgreSQL + ClickHouse)
- [x] API contracts (REST + gRPC)
- [x] Ingestion pipeline (Kafka + Temporal + outbox)
- [x] Security model (OpenFGA + Argon2id)
- [x] Observability SLIs mapped to NFRs
- [x] Infrastructure and deployment path
- [x] Docker Compose profiles
- [x] 8 ADRs
- [ ] Architecture review by senior engineer (pending)
- [ ] Threat model sign-off (pending)

## Next: Phase 1

Scaffold Nx monorepo, implement 6 Go services, Next.js dashboard, Node SDK, Docker Compose `services` profile.