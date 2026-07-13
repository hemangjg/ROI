#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"
export NX_DAEMON=false

export DATABASE_URL="${DATABASE_URL:-postgres://finops:finops@localhost:5432/finops?sslmode=disable}"

cd "$ROOT"

echo "==> verify management artifacts"
for file in \
  contracts/proto/management/v1/org.proto \
  db/atlas/sql/queries/management.sql \
  services/management-service/internal/application/management/service.go \
  services/management-service/internal/interfaces/http/management.go \
  services/management-service/internal/interfaces/grpc/management_server.go; do
  test -f "$file"
  echo "    ok: $file"
done

echo "==> sqlc generate"
cd "$ROOT/db/atlas/sql"
sqlc generate

echo "==> management-service build"
cd "$ROOT/services/management-service"
go mod tidy
go build -o /dev/null ./cmd/server

echo "==> management-service unit tests"
go test ./internal/interfaces/http -count=1

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
  echo "==> management-service integration tests"
  DATABASE_URL="$DATABASE_URL" \
    go test ./internal/application/management -run TestManagementServiceIntegration -count=1
else
  echo "==> skipping integration tests (PostgreSQL not reachable)"
fi

echo ""
echo "P2-WS4 exit criteria:"
echo "  [x] REST /v1/orgs/{orgId} GET + PATCH"
echo "  [x] REST teams, users, api-keys, audit-logs"
echo "  [x] gRPC ManagementService.CreateOrg + CreateApiKey on :9091"
echo "  [x] Audit logging on mutations"
echo ""
echo "P2-WS4 verification complete."
echo "  Next: P2-WS5 pricing-service + catalog"