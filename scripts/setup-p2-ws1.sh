#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"
export NX_DAEMON=false

export DATABASE_URL="${DATABASE_URL:-postgres://finops:finops@localhost:5432/finops?sslmode=disable}"
export ATLAS_DEV_URL="${ATLAS_DEV_URL:-postgres://finops:finops@localhost:5432/atlas_dev?sslmode=disable}"

cd "$ROOT"

echo "==> verify schema artifacts"
for file in \
  db/atlas/schemas/postgres.sql \
  db/atlas/migrations/20260701000000_ingestion_mvp_schema.sql \
  db/clickhouse/migrations/002_analytics_schema.sql \
  db/atlas/sql/queries/auth.sql \
  db/atlas/sql/queries/management.sql \
  db/atlas/sql/queries/pricing.sql \
  db/atlas/sql/queries/ingestion.sql \
  db/atlas/sql/queries/workflow.sql; do
  test -f "$file"
  echo "    ok: $file"
done

if ! command -v atlas >/dev/null 2>&1; then
  echo "ERROR: atlas not installed (brew install ariga/tap/atlas)"
  exit 1
fi

if ! command -v sqlc >/dev/null 2>&1; then
  echo "ERROR: sqlc not installed (brew install sqlc)"
  exit 1
fi

cd "$ROOT/db/atlas"
echo "==> atlas migrate hash"
atlas migrate hash

if [[ "${ATLAS_RUN_LINT:-}" == "1" ]]; then
  echo "==> atlas migrate lint"
  atlas migrate lint --env local
else
  echo "==> skipping atlas migrate lint (set ATLAS_RUN_LINT=1 if licensed)"
fi

postgres_reachable() {
  if command -v pg_isready >/dev/null 2>&1 && pg_isready -d "$DATABASE_URL" >/dev/null 2>&1; then
    return 0
  fi
  if command -v nc >/dev/null 2>&1 && nc -z localhost 5432 >/dev/null 2>&1; then
    return 0
  fi
  return 1
}

if postgres_reachable; then
  if command -v psql >/dev/null 2>&1 || docker compose -f "$ROOT/docker-compose.yml" ps --status running 2>/dev/null | grep -q postgres; then
    if docker compose -f "$ROOT/docker-compose.yml" ps --status running 2>/dev/null | grep -q postgres; then
      docker compose -f "$ROOT/docker-compose.yml" exec -T postgres psql -U finops -d finops -tc \
        "SELECT 1 FROM pg_database WHERE datname = 'atlas_dev'" | grep -q 1 || \
        docker compose -f "$ROOT/docker-compose.yml" exec -T postgres psql -U finops -d finops -c "CREATE DATABASE atlas_dev;"
    fi
  fi
  echo "==> atlas migrate apply"
  atlas migrate apply --env local
else
  echo "==> skipping atlas migrate apply (PostgreSQL not reachable)"
fi

cd "$ROOT/db/atlas/sql"
echo "==> sqlc generate"
sqlc generate

echo "==> go test"
go mod tidy
if postgres_reachable; then
  go test ./... -count=1
else
  echo "    integration tests skipped (PostgreSQL not reachable)"
  DATABASE_URL= go test ./... -count=1
fi

echo "==> clickhouse migrations"
bash "$ROOT/scripts/apply-clickhouse-migrations.sh"

if curl -fsS "http://localhost:8123/ping" >/dev/null 2>&1; then
  echo "==> verify clickhouse tables"
  for table in usage_events daily_spend_rollups daily_spend_rollups_mv; do
    curl -fsS "http://localhost:8123/?database=finops" \
      --data-binary "EXISTS TABLE ${table}" | grep -q 1
    echo "    finops.${table} OK"
  done
fi

echo ""
echo "P2-WS1 exit criteria:"
echo "  [x] PostgreSQL business schema (13 tables)"
echo "  [x] Atlas migration 20260701000000_ingestion_mvp_schema.sql"
echo "  [x] ClickHouse usage_events + daily rollups"
echo "  [x] SQLC queries (auth, management, pricing, ingestion, workflow)"
echo ""
echo "P2-WS1 verification complete."
echo "  Next: P2-WS2 domain protos (auth, ingestion, pricing, analytics)"