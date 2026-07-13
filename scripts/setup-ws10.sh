#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"

echo "==> pnpm install"
pnpm install

echo "==> build api-schemas"
pnpm --filter @ai-finops/api-schemas build

echo "==> test sdk-core"
pnpm --filter @ai-finops/sdk-core test

echo "==> build sdk-core"
pnpm --filter @ai-finops/sdk-core build

echo "==> test sdk-typescript"
pnpm --filter @ai-finops/sdk test

echo "==> build sdk-typescript"
pnpm --filter @ai-finops/sdk build

echo "==> verify dist artifacts"
test -f packages/sdk-typescript/dist/index.js
test -f packages/sdk-typescript/dist/index.cjs
test -f packages/sdk-typescript/dist/index.d.ts

echo "==> nx build targets"
pnpm exec nx run sdk-typescript:build

echo "WS10 verification complete."
echo "  Package: @ai-finops/sdk"
echo "  Import:  import { FinOpsClient } from '@ai-finops/sdk';"
echo "  API:     POST /events via client.ingest(...)"