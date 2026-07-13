# ADR-002: Nx Monorepo with Polyglot Layout

## Status

Accepted

## Context

The platform spans Go microservices, a Next.js dashboard, and multi-language SDKs. Shared API contracts (OpenAPI, Protobuf) must stay in sync across all consumers.

## Decision

Use **Nx + pnpm workspaces** with this layout:
- `services/` — 6 Go microservices (independent `go.mod` per service)
- `apps/web/` — Next.js dashboard
- `packages/` — SDKs and generated API schemas
- `api/` — OpenAPI and Buf proto definitions (source of truth)
- `gen/go/` — Buf-generated Go stubs (gitignored)

Nx orchestrates build, test, lint across all projects. `nx affected` runs only changed projects in CI.

## Consequences

**Positive:**
- Single repo for API contract changes with atomic PRs
- Nx caching speeds CI
- Clear boundary: services communicate via gRPC only

**Negative:**
- Polyglot repo complexity (Go + TypeScript tooling)
- Developers need both Go and Node.js environments

## Alternatives Considered

| Alternative | Rejected because |
|-------------|-----------------|
| Turborepo | Nx has stronger affected-graph for polyglot |
| Multi-repo | API contract sync across repos is painful |
| Go monolith | Cannot scale ingest independently from analytics |