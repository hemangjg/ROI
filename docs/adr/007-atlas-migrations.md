# ADR-007: Atlas for PostgreSQL Migrations

## Status

Accepted

## Context

Database schema evolves frequently (new providers, governance tables, budget features). Migrations must be versioned, reviewable, and detect drift between environments.

## Decision

Use **Atlas** for PostgreSQL schema management:
- Declarative schema in `db/atlas/schemas/postgres.sql`
- Versioned migrations in `db/atlas/migrations/`
- `atlas migrate diff` generates migrations from schema changes
- `atlas schema inspect` detects drift in staging/production
- Environment configs in `db/atlas/atlas.hcl`

ClickHouse migrations use versioned SQL in `db/clickhouse/migrations/` (Atlas does not support ClickHouse natively).

## Consequences

**Positive:**
- Declarative schema is single source of truth
- Automatic diff reduces manual migration writing
- Drift detection prevents environment inconsistencies
- Multi-environment HCL configs (local, staging, production)

**Negative:**
- Atlas is newer than golang-migrate (smaller community)
- ClickHouse requires separate migration tooling
- Team must learn Atlas CLI

## Alternatives Considered

| Alternative | Rejected because |
|-------------|-----------------|
| golang-migrate | Manual migration writing; no declarative diff; no drift detection |
| Flyway | Java-centric; less Go ecosystem integration |
| Prisma Migrate | Backend is Go, not Node |
| Hand-written SQL only | No drift detection; error-prone |