# Database Foundation

Phase 1b (Ingestion MVP) adds the business schema. Phase 1 provided toolchain wiring only.

## PostgreSQL (Atlas + SQLC + pgx)

```
db/atlas/
├── atlas.hcl
├── schemas/postgres.sql      # declarative schema (13 tables)
├── migrations/
│   ├── 20260628000000_baseline.sql
│   └── 20260701000000_ingestion_mvp_schema.sql
└── sql/
    ├── sqlc.yaml
    ├── queries/              # auth, management, pricing, ingestion, workflow
    └── gen/                  # sqlc output (committed)
```

### Tables

| Table | Owner service | Purpose |
|-------|---------------|---------|
| `organizations` | management | Tenant root |
| `users` | identity | Global accounts |
| `org_memberships` | management | User ↔ org roles |
| `teams` | management | Soft-deletable teams |
| `team_members` | management | User ↔ team roles |
| `api_keys` | management | Ingest authentication |
| `refresh_tokens` | identity | Rotating dashboard sessions |
| `providers` | pricing | Provider catalog |
| `models` | pricing | Model catalog |
| `model_pricing` | pricing | Versioned pricing rows |
| `outbox_events` | ingestion | Transactional outbox |
| `idempotency_keys` | workflow | Dedupe processed events |
| `audit_logs` | management | Append-only audit trail |

### Commands

```bash
# Verify P2-WS1 (schema + sqlc + optional apply)
bash scripts/setup-p2-ws1.sh

# Apply migrations (requires PostgreSQL)
export DATABASE_URL=postgres://finops:finops@localhost:5432/finops?sslmode=disable
atlas migrate apply --env local --dir file://db/atlas/migrations

# Generate SQLC
pnpm sqlc:generate
```

## ClickHouse

```
db/clickhouse/migrations/
├── 001_baseline.sql
└── 002_analytics_schema.sql   # usage_events + daily_spend_rollups
```

```bash
bash scripts/apply-clickhouse-migrations.sh
```

Analytics Service reads `daily_spend_rollups`; only Workflow Service writes `usage_events`.