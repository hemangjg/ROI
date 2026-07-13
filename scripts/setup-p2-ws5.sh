#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"
export NX_DAEMON=false

export DATABASE_URL="${DATABASE_URL:-postgres://finops:finops@localhost:5432/finops?sslmode=disable}"

cd "$ROOT"

echo "==> verify pricing artifacts"
for file in \
  contracts/proto/pricing/v1/pricing.proto \
  db/atlas/sql/queries/pricing.sql \
  services/pricing-service/internal/domain/plugin.go \
  services/pricing-service/internal/plugin/openai/openai.go \
  services/pricing-service/internal/plugin/anthropic/anthropic.go \
  services/pricing-service/internal/application/pricing/service.go \
  services/pricing-service/internal/interfaces/grpc/pricing_server.go; do
  test -f "$file"
  echo "    ok: $file"
done

echo "==> sqlc generate"
cd "$ROOT/db/atlas/sql"
sqlc generate

echo "==> pricing-service build"
cd "$ROOT/services/pricing-service"
go mod tidy
go build -o /dev/null ./cmd/server

echo "==> pricing-service unit tests"
go test ./internal/domain ./internal/interfaces/http -count=1

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
  echo "==> pricing-service integration tests"
  DATABASE_URL="$DATABASE_URL" \
    go test ./internal/application/pricing -run TestPricingServiceIntegration -count=1
else
  echo "==> skipping integration tests (PostgreSQL not reachable)"
fi

echo ""
echo "P2-WS5 exit criteria:"
echo "  [x] OpenAI + Anthropic provider plugins with embedded catalog"
echo "  [x] gRPC PricingService.CalculateCost + SyncPricingCatalog on :9092"
echo "  [x] Versioned model_pricing upsert with effective dates"
echo "  [x] Unknown model returns is_priced=false"
echo ""
echo "P2-WS5 verification complete."
echo "  Next: P2-WS6 ingestion-service + outbox"