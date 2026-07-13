# ADR-004: Kafka + Temporal + Transactional Outbox

## Status

Accepted

## Context

Usage event ingestion must return in < 100ms while guaranteeing durability. Direct Kafka publish from the API handler risks message loss if Kafka is unavailable. Cost calculation involves multiple steps with retry requirements.

## Decision

1. **Transactional outbox:** Ingestion Service writes `outbox_events` in PostgreSQL TX, returns 202
2. **Outbox relay:** Polls unpublished rows, publishes to Kafka `usage.events.raw`
3. **Kafka:** Durable event stream with 7-day retention
4. **Temporal:** `ProcessUsageEvent` workflow orchestrates dedupe → price → write → record
5. **Dead-letter queue:** `usage.events.dlq` for exhausted retries

## Consequences

**Positive:**
- Ingest API never blocks on Kafka or ClickHouse availability
- Temporal provides durable retries with visibility
- Kafka enables replay and multiple consumers (future analytics pipelines)
- Outbox guarantees at-least-once delivery

**Negative:**
- Three infrastructure components (PostgreSQL outbox, Kafka, Temporal)
- Eventual consistency (target < 30s processing lag)
- Operational learning curve for Temporal

## Alternatives Considered

| Alternative | Rejected because |
|-------------|-----------------|
| BullMQ / Redis queue | No durable replay; weaker persistence guarantees |
| Synchronous write to ClickHouse | Violates p99 < 100ms ingest target |
| Direct Kafka (no outbox) | Message loss if Kafka down during publish |
| Celery | Python-centric; team standard is Go |