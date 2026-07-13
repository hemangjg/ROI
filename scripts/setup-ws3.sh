#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

if ! command -v atlas >/dev/null 2>&1; then
  echo "Installing atlas..."
  brew install ariga/tap/atlas
fi

if ! command -v sqlc >/dev/null 2>&1; then
  echo "Installing sqlc..."
  brew install sqlc
fi

export DATABASE_URL="${DATABASE_URL:-postgres://finops:finops@localhost:5432/finops?sslmode=disable}"
export ATLAS_DEV_URL="${ATLAS_DEV_URL:-postgres://finops:finops@localhost:5432/atlas_dev?sslmode=disable}"

postgres_reachable() {
  if command -v pg_isready >/dev/null 2>&1 && pg_isready -d "$DATABASE_URL" >/dev/null 2>&1; then
    return 0
  fi
  if command -v nc >/dev/null 2>&1 && nc -z localhost 5432 >/dev/null 2>&1; then
    return 0
  fi
  return 1
}

ensure_atlas_dev_db() {
  if docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^finops-pg$'; then
    if ! docker exec finops-pg psql -U finops -d finops -tc \
      "SELECT 1 FROM pg_database WHERE datname = 'atlas_dev'" | grep -q 1; then
      docker exec finops-pg psql -U finops -d finops -c "CREATE DATABASE atlas_dev;"
    fi
    return 0
  fi

  if command -v psql >/dev/null 2>&1; then
    psql "$DATABASE_URL" -tc \
      "SELECT 1 FROM pg_database WHERE datname = 'atlas_dev'" | grep -q 1 || \
      psql "$DATABASE_URL" -c "CREATE DATABASE atlas_dev;"
    return 0
  fi

  return 1
}

cd "$ROOT/db/atlas"

echo "==> atlas migrate hash"
atlas migrate hash

if [[ "${ATLAS_RUN_LINT:-}" == "1" ]]; then
  if postgres_reachable; then
    ensure_atlas_dev_db || echo "    Could not create atlas_dev; lint may fail."
    echo "==> atlas migrate lint (Atlas Pro; uses atlas_dev on localhost:5432)"
    atlas migrate lint --env local
  else
    echo "==> skipping atlas migrate lint (PostgreSQL not reachable at localhost:5432)"
  fi
else
  echo "==> skipping atlas migrate lint (Atlas Pro since v0.38; set ATLAS_RUN_LINT=1 if licensed)"
fi

if postgres_reachable; then
  echo "==> atlas migrate apply"
  atlas migrate apply --env local
else
  echo "==> skipping atlas migrate apply (PostgreSQL not reachable at localhost:5432)"
  echo "    Start PostgreSQL or set DATABASE_URL to apply migrations."
fi

cd "$ROOT/db/atlas/sql"

echo "==> sqlc generate"
sqlc generate

echo "==> go test (integration test skips without DATABASE_URL + running Postgres)"
go mod tidy
go test ./... -count=1

echo "WS3 verification complete."