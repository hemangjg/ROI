#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"

echo "==> pnpm install"
pnpm install

echo "==> verify dev scripts"
test -x scripts/dev.sh
test -x scripts/dev-preflight.sh
test -f .github/workflows/ci.yml

echo "==> verify pnpm dev wiring"
pnpm exec nx show project workspace --json | grep -q '"dev"'

echo "==> CI local checks"
SKIP_WEB_BUILD=1 bash scripts/ci-local.sh

if docker info >/dev/null 2>&1; then
  echo "==> dev preflight (docker core only)"
  bash scripts/dev-preflight.sh
else
  echo "Docker not running — skipped dev preflight"
fi

echo "WS11 verification complete."
echo "  Native dev: pnpm dev"
echo "  Docker dev: pnpm dev:docker  (or: make dev)"
echo "  CI local:   bash scripts/ci-local.sh"