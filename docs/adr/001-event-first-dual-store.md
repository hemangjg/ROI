# ADR-001: Event-First Dual-Store Data Model

## Status

Accepted

## Context

AI FinOps must ingest millions of usage events per day per organization while serving sub-second dashboard queries. A single PostgreSQL database cannot efficiently handle both write-heavy telemetry and read-heavy analytics at enterprise scale.

## Decision

1. **Usage events are immutable** — append-only in ClickHouse; never UPDATE cost after write
2. **PostgreSQL stores business data** — orgs, users, teams, API keys, pricing catalog, audit logs, outbox
3. **ClickHouse stores telemetry** — usage_events and materialized daily_spend_rollups
4. **Rollups are derived** — rebuildable from raw events via RecomputeRollups workflow
5. **Pricing changes never retroactively alter events** — price lookup uses `occurred_at` effective date

## Consequences

**Positive:**
- ClickHouse handles 500M+ events/day platform-wide
- Dashboard queries scan rollups only (< 500ms p95)
- Clear separation of concerns per service

**Negative:**
- Operational complexity of two databases
- Eventual consistency between ingest and dashboard (target < 30s)
- ClickHouse migrations managed separately from Atlas

## Alternatives Considered

| Alternative | Rejected because |
|-------------|-----------------|
| Single PostgreSQL | Won't scale past ~10M events/day without expensive partitioning |
| PostgreSQL + TimescaleDB | Less proven at 500M events/day; ClickHouse better for analytics |
| Event sourcing only | Over-engineered for Phase 1; immutability without full event sourcing |