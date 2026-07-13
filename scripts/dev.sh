#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}/usr/local/bin:${HOME}/go/bin"

cd "$ROOT"

bash "$ROOT/scripts/dev-preflight.sh"

echo "==> starting native Go services + Next.js dashboard"
exec pnpm exec concurrently \
  --kill-others-on-fail \
  --prefix "[{name}]" \
  --names "identity,ingestion,management,analytics,pricing,workflow,web" \
  --prefix-colors "blue,green,yellow,cyan,magenta,red,white" \
  "nx run identity-service:serve" \
  "nx run ingestion-service:serve" \
  "nx run management-service:serve" \
  "nx run analytics-service:serve" \
  "nx run pricing-service:serve" \
  "nx run workflow-service:serve" \
  "nx run web:serve"