# AI FinOps — verification status

**Last verified:** 2026-07-13  
**Repo path:** `/Users/hemangjg/Hemangs stuff/Projects/ai-finops`

## Phase 1 + 1b engineering sign-off

| Check | Result |
|-------|--------|
| `pnpm format:check` + web ESLint | Pass |
| `bash scripts/lint-go.sh` (golangci-lint) | Pass |
| `bash scripts/ci-local.sh` | Pass |
| Docker core + services stack | Pass (local) |
| `scripts/setup-p2-ws1.sh` … `setup-p2-ws9.sh` | Pass with Postgres/ClickHouse |
| `bash scripts/e2e-smoke.sh` | **Pass** — register → API key → ingest → spend summary |

## E2E proof (sample)

```
register → create API key → POST /v1/events → workflow prices event →
GET /v1/orgs/{id}/spend/summary shows non-zero today_usd / mtd_usd
```

## Fixes landed during completion pass

1. Prettier formatting on 8 files  
2. golangci-lint errcheck fixes (`packages/config`, analytics ClickHouse rows, outbox/ingestion closes)  
3. ClickHouse insert panic on nil `*uuid.UUID` (team/user metadata)  
4. Envoy spend route regex was a **full-path** match and never hit `/spend/summary`  
5. Kafka topic auto-create + init ensures `usage.events.raw`  
6. Outbox writer `AllowAutoTopicCreation`  
7. Honest Phase 1b / path docs; `scripts/e2e-smoke.sh` added  
8. Workflow consumer may need restart after Kafka restarts (0 partitions) — handled in `init-local.sh`

## How to re-verify

```bash
cd "/Users/hemangjg/Hemangs stuff/Projects/ai-finops"
export PATH="$(go env GOPATH)/bin:$PATH"
bash scripts/quality-gate.sh
# Live pipeline:
WITH_E2E=1 bash scripts/quality-gate.sh
```

See also [runbook.md](./runbook.md) and root `AGENTS.md`.

## Base hardening (post Phase 1b)

- `packages/events.EnsureTopic` on outbox + workflow boot
- Kafka reader `WatchPartitionChanges`
- Unit tests: workflow UUID helpers, OpenAI/Anthropic catalogs, events validation
- `scripts/quality-gate.sh` single entry for offline (+ optional e2e)
- `AGENTS.md` + `docs/runbook.md`

## Phase 2 — budgets

Shipped (Sprint A):

- `budgets` + `budget_alerts` tables (Atlas migration `20260713120000_budgets.sql`)
- Management API: `PUT/GET /v1/orgs/{id}/budgets`, `POST .../budgets/evaluate`, `GET/POST .../budget-alerts`
- Soft threshold unit tests (Go + web helpers)
- Dashboard **Budgets**: live MTD analytics by-team spend, progress bars + soft marker, **Evaluate live spend**
- Overview open-alert card + sidebar Budgets badge
- `scripts/e2e-budgets.sh` (register → team → budget → soft alert → ack)
- OpenAPI updated

Still later: email/Slack delivery, hard limits, forecasting, model allowlist.
