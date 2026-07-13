#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"

services=(
  identity-service
  ingestion-service
  management-service
  analytics-service
  pricing-service
  workflow-service
)

for svc in "${services[@]}"; do
  echo "==> testing services/$svc"
  (cd "services/$svc" && go test ./... -count=1)
done

echo "All service tests passed."