#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"

echo "==> pnpm install"
pnpm install

echo "==> lint web"
pnpm --filter @ai-finops/web lint

echo "==> build web"
pnpm --filter @ai-finops/web build

if docker info >/dev/null 2>&1; then
  echo "==> docker build web image"
  docker compose --profile web build web

  echo "==> start web container"
  docker compose --profile web up -d web

  echo "==> waiting for dashboard"
  deadline=$((SECONDS + 120))
  ready=false
  while (( SECONDS < deadline )); do
    if curl -fsS http://localhost:3000/login 2>/dev/null | grep -q "AI FinOps"; then
      ready=true
      break
    fi
    sleep 3
  done

  if [[ "$ready" != "true" ]]; then
    echo "ERROR: dashboard did not become ready"
    docker compose logs web --tail 50
    exit 1
  fi

  echo "==> route checks"
  curl -fsS http://localhost:3000/login | grep -q "AI FinOps"
  echo "    /login OK"
else
  echo "Docker not running — skipped container verification"
fi

echo "WS9 verification complete."
echo "  Local dev: cd apps/web && pnpm dev"
echo "  Dashboard: http://localhost:3000"