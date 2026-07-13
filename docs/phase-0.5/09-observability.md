# 09 — Observability

## Stack

| Signal | Collection | Storage | Visualization |
|--------|-----------|---------|---------------|
| Traces | OpenTelemetry SDK | Jaeger | Jaeger UI / Grafana |
| Metrics | OTel → Prometheus exporter | Prometheus | Grafana |
| Logs | slog → OTel Logs | Loki | Grafana |
| Alerts | Prometheus rules | Alertmanager | Slack / PagerDuty |

All Go services and Next.js instrumented from day one.

---

## OpenTelemetry Instrumentation

### Go services

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)
```

- HTTP: `otelhttp.NewHandler` on Chi routes
- gRPC: `otelgrpc.UnaryServerInterceptor` / `UnaryClientInterceptor`
- Custom spans: `ProcessUsageEvent`, `CalculateCost`, `ClickHouseInsert`

### Trace propagation

W3C Trace Context headers propagated: SDK → Envoy → Ingestion → Kafka (baggage) → Workflow → Pricing → ClickHouse.

### Resource attributes

```
service.name = ingestion-service
service.version = 1.0.0
deployment.environment = production
org.id = (when available, not on ingest to avoid cardinality explosion)
```

---

## Key Metrics (Prometheus)

### Ingestion Service

| Metric | Type | Labels |
|--------|------|--------|
| `ingest_requests_total` | Counter | status, org_id (sampled) |
| `ingest_latency_seconds` | Histogram | — |
| `ingest_events_accepted_total` | Counter | org_id (sampled) |
| `ingest_rate_limit_hits_total` | Counter | key_id |

### Workflow Service

| Metric | Type | Labels |
|--------|------|--------|
| `kafka_consumer_lag` | Gauge | topic, partition |
| `workflow_duration_seconds` | Histogram | workflow_name |
| `workflow_failures_total` | Counter | workflow_name, activity |
| `unpriced_events_total` | Counter | provider, model |

### Analytics Service

| Metric | Type | Labels |
|--------|------|--------|
| `query_duration_seconds` | Histogram | endpoint |
| `query_cache_hits_total` | Counter | endpoint |

### Pricing Service

| Metric | Type | Labels |
|--------|------|--------|
| `pricing_catalog_age_days` | Gauge | provider |
| `cost_calculation_duration_seconds` | Histogram | — |

---

## SLI / SLO Mapping

| NFR | SLI | SLO target |
|-----|-----|------------|
| NFR-PERF-001 | `ingest_latency_seconds` p99 | < 100ms |
| NFR-PERF-003 | `kafka_consumer_lag` p99 | < 30s equivalent |
| NFR-PERF-005 | `query_duration_seconds` p95 | < 500ms |
| NFR-AVAIL-003 | `up{job="ingestion-service"}` | 99.95% |
| FR-COST accuracy | `unpriced_events_total` rate | < 0.5% of total |

---

## Grafana Dashboards

| Dashboard | Audience | Panels |
|-----------|----------|--------|
| Ingest Health | Engineering | RPS, latency, error rate, rate limits |
| Pipeline Lag | Engineering | Kafka lag, Temporal queue depth, DLQ size |
| Cost Accuracy | FinOps ops | Unpriced events, pricing catalog age |
| Service Overview | SRE | Per-service up, CPU, memory, request rate |

---

## Alerting Rules (Alertmanager)

| Alert | Condition | Severity |
|-------|-----------|----------|
| IngestHighErrorRate | 5xx > 1% for 5min | P1 |
| KafkaLagHigh | lag > 10000 for 5min | P1 |
| PricingCatalogStale | age > 7 days | P2 |
| UnpricedEventsSpike | rate > 100/hr | P2 |
| ClickHouseDown | up == 0 for 1min | P1 |

---

## Structured Logging (slog)

```json
{
  "time": "2026-06-28T10:00:00Z",
  "level": "INFO",
  "msg": "event ingested",
  "trace_id": "abc123",
  "org_id": "uuid",
  "event_id": "uuid",
  "service": "ingestion-service"
}
```

- No PII in logs (no prompt content, no API key values)
- `trace_id` correlates with Jaeger
- Shipped to Loki via Promtail or OTel Collector logs pipeline