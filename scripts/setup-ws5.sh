#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"

echo "==> go work sync"
go work sync

services=(
  identity-service
  ingestion-service
  management-service
  analytics-service
  pricing-service
  workflow-service
)

for svc in "${services[@]}"; do
  echo "==> go mod tidy services/$svc"
  (cd "services/$svc" && go mod tidy)
done

echo "==> build all services"
mkdir -p dist
for svc in "${services[@]}"; do
  echo "    building $svc"
  (cd "services/$svc" && go build -o "$ROOT/dist/$svc" ./cmd/server)
done

echo "==> test all services"
bash "$ROOT/scripts/test-services.sh"

echo "WS5 verification complete."