# 03 — API Contracts

## API Strategy

| Audience | Protocol | Spec location |
|----------|----------|---------------|
| SDKs, external integrators | REST | `api/openapi/v1.yaml` |
| Service-to-service | gRPC + Protobuf | `api/proto/` (Buf module) |
| Envoy ext_authz | gRPC | `auth.v1.AuthService.ValidateToken` |

**ADR-005:** REST external, gRPC internal. **ADR-008:** Buf manages proto lifecycle.

---

## External REST API

Base URL: `https://api.{domain}/v1`

### Authentication

| Path type | Header |
|-----------|--------|
| Ingest | `Authorization: Bearer aif_...` |
| Dashboard/Admin | `Authorization: Bearer <jwt>` |

### Ingestion Service endpoints

#### POST /events

Ingest single usage event. Returns `202 Accepted`.

Request:
```json
{
  "idempotency_key": "req_abc123",
  "provider": "openai",
  "model": "gpt-4o",
  "input_tokens": 1200,
  "output_tokens": 340,
  "cache_read_tokens": 0,
  "cache_write_tokens": 0,
  "occurred_at": "2026-06-28T10:00:00Z",
  "metadata": {
    "team_id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "project": "checkout-bot",
    "feature": "summarize",
    "customer": "acme-corp"
  }
}
```

Response `202`:
```json
{ "event_id": "uuid", "status": "queued" }
```

#### POST /events/batch

Max 500 events. Partial success supported.

#### GET /events/{eventId}

Returns processing status and cost when complete.

---

### identity-service endpoints

| Method | Path | Body | Response |
|--------|------|------|----------|
| POST | `/auth/register` | email, password, name, org_name | JWT pair + org_id |
| POST | `/auth/login` | email, password | JWT pair |
| POST | `/auth/refresh` | refresh_token | new JWT pair |
| POST | `/auth/logout` | refresh_token | 204 |

---

### Management Service endpoints

| Method | Path | Role | Purpose |
|--------|------|------|---------|
| GET | `/orgs/{orgId}` | viewer | Get org details |
| PATCH | `/orgs/{orgId}` | admin | Update org settings |
| POST | `/orgs/{orgId}/teams` | admin | Create team |
| GET | `/orgs/{orgId}/teams` | viewer | List teams |
| DELETE | `/orgs/{orgId}/teams/{teamId}` | admin | Soft-delete team |
| POST | `/orgs/{orgId}/users/invite` | admin | Invite user |
| GET | `/orgs/{orgId}/users` | admin | List members |
| POST | `/orgs/{orgId}/api-keys` | admin | Create key (plaintext once) |
| GET | `/orgs/{orgId}/api-keys` | admin | List keys (prefix only) |
| DELETE | `/orgs/{orgId}/api-keys/{keyId}` | admin | Revoke key |
| GET | `/orgs/{orgId}/audit-logs` | admin | Paginated audit trail |

---

### Analytics Service endpoints

| Method | Path | Query params | Response |
|--------|------|-------------|----------|
| GET | `/orgs/{orgId}/spend/summary` | from, to, team_id | today, mtd, change_pct |
| GET | `/orgs/{orgId}/spend/timeseries` | from, to, granularity, team_id | daily data points |
| GET | `/orgs/{orgId}/spend/by-provider` | from, to, team_id | provider breakdown |
| GET | `/orgs/{orgId}/spend/by-team` | from, to | team table |

All dates: ISO 8601. Default range: last 30 days.

---

## Internal gRPC Services

Generated via `buf generate` into `gen/go/`.

| Package | Service | Key RPCs |
|---------|---------|----------|
| `auth.v1` | AuthService | ValidateToken, ValidateApiKey |
| `ingestion.v1` | IngestionService | IngestEvent |
| `management.v1` | ManagementService | CreateOrg, CreateApiKey |
| `analytics.v1` | AnalyticsService | GetSpendSummary, GetSpendTimeSeries |
| `pricing.v1` | PricingService | CalculateCost, SyncPricingCatalog |
| `workflow.v1` | WorkflowService | ProcessUsageEvent |

Proto files: `api/proto/{auth,ingestion,management,analytics,pricing,workflow}/v1/*.proto`

### Buf commands

```bash
cd api/proto
buf lint
buf breaking --against '.git#branch=main'
buf generate
```

---

## Error Format

All REST errors use RFC 7807 Problem Details:

```json
{
  "type": "https://api.ai-finops.io/errors/validation",
  "title": "Validation Error",
  "status": 400,
  "detail": "model is required",
  "instance": "/v1/events"
}
```

| Status | Usage |
|--------|-------|
| 400 | Validation error |
| 401 | Missing/invalid auth |
| 403 | OpenFGA denied |
| 404 | Resource not found |
| 429 | Rate limit exceeded |
| 202 | Event accepted (ingest) |
| 500 | Internal error |

---

## Rate Limits

| Endpoint | Limit |
|----------|-------|
| POST /events | 1,000 events/min per API key |
| POST /events/batch | Counts as N events toward limit |
| Dashboard APIs | 100 req/min per user |
| Auth endpoints | 10 req/min per IP |

`429` response includes `Retry-After` header.

---

## Versioning

- REST: URL prefix `/v1/`; breaking changes require `/v2/`
- gRPC: Buf breaking change detection in CI
- SDK: SemVer aligned with API version