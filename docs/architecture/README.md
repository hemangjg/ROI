# Architecture Diagrams (C4)

Living C4 views for the AI FinOps platform. Phase 1 foundation implements shells and routing only; business APIs land in Phase 2.

| Diagram | File | Level |
|---------|------|-------|
| Context | [01-context.md](./01-context.md) | Who uses the system and what it connects to |
| Container | [02-container.md](./02-container.md) | Deployable services, data stores, messaging |
| Component | [03-component.md](./03-component.md) | Phase 1 monorepo layout and service internals |

## Related docs

- [Phase 0.5 system architecture](../phase-0.5/01-system-architecture.md) — full design narrative
- [ADRs](../adr/) — locked decisions
- [OpenAPI contract](../../api/openapi/v1.yaml) — external REST API
- [Proto contracts](../../contracts/proto/) — internal gRPC (bootstrap in Phase 1)

## Phase 1 vs Phase 2

| Area | Phase 1 (foundation) | Phase 2+ |
|------|----------------------|----------|
| Services | Health/readiness shells, OTel, SQLC wiring | Domain handlers, gRPC servers |
| Database | Empty Atlas baseline | Full entity model + migrations |
| Protos | `bootstrap/health` only | Auth, ingestion, pricing, analytics |
| Gateway | Envoy routing + TLS | ext_authz JWT validation |
| Dashboard | UI shell + mock auth | Live API integration |