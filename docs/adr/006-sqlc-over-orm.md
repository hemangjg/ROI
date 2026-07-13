# ADR-006: SQLC over ORM

## Status

Accepted

## Context

PostgreSQL queries must be tenant-scoped, performant, and type-safe. ORMs (GORM, Ent) hide SQL, make tenant isolation enforcement harder, and generate inefficient queries for analytics-adjacent lookups.

## Decision

Use **SQLC** to generate type-safe Go code from raw SQL queries in `sql/queries/`. Schema managed by **Atlas** (see ADR-007). Each Go service has its own SQLC config scoped to its queries.

## Consequences

**Positive:**
- Explicit SQL — every query reviewed in PR
- Compile-time type safety on query results
- No ORM magic; predictable query plans
- Natural enforcement of `WHERE org_id = $1` in every query

**Negative:**
- Manual SQL writing (not auto-generated from structs)
- SQLC regeneration step in dev workflow
- No automatic relationship loading (explicit joins)

## Alternatives Considered

| Alternative | Rejected because |
|-------------|-----------------|
| GORM | Hidden queries; tenant leak risk; performance unpredictability |
| Ent | Heavy codegen; steeper learning curve |
| Prisma (Node) | Backend is Go, not Node |
| Raw pgx | No compile-time type safety on scan |