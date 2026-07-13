#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"
export NX_DAEMON=false

cd "$ROOT"

echo "==> install dependencies"
pnpm install
go work sync

echo "==> lint go"
bash scripts/lint-go.sh

echo "==> lint typescript"
pnpm format:check
pnpm --filter @ai-finops/web lint

echo "==> proto lint + generate"
if command -v buf >/dev/null 2>&1; then
  pnpm proto:lint
  if [[ -d packages/proto/gen/go && -d packages/proto/gen/ts ]]; then
    echo "    proto gen present; skipping buf generate (run pnpm proto:generate to refresh)"
  else
    pnpm proto:generate
  fi
else
  echo "buf not installed; skipping proto checks"
fi

echo "==> atlas migrate hash"
if command -v atlas >/dev/null 2>&1; then
  (
    cd db/atlas
    atlas migrate hash
    if [[ "${ATLAS_RUN_LINT:-}" == "1" ]]; then
      export DATABASE_URL="${DATABASE_URL:-postgres://finops:finops@localhost:5432/finops?sslmode=disable}"
      export ATLAS_DEV_URL="${ATLAS_DEV_URL:-postgres://finops:finops@localhost:5432/atlas_dev?sslmode=disable}"
      atlas migrate lint --env local
    else
      echo "    skipping atlas migrate lint (Atlas Pro; set ATLAS_RUN_LINT=1 to enable)"
    fi
  )
else
  echo "atlas not installed; skipping atlas checks"
fi

echo "==> codegen drift"
bash scripts/check-codegen-drift.sh

echo "==> test"
pnpm test:go
bash scripts/test-services.sh
pnpm exec nx run-many -t test --projects=web,sdk-core,sdk-typescript

echo "==> build (Go services + SDK packages)"
pnpm exec nx run-many -t build --exclude=web,workspace

if [[ "${SKIP_WEB_BUILD:-}" == "1" ]]; then
  echo "==> skipping web production build (SKIP_WEB_BUILD=1)"
else
  echo "==> web production build"
  bash scripts/build-web.sh
fi

echo "CI local verification complete."