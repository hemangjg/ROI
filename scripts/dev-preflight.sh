#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"

if ! docker info >/dev/null 2>&1; then
  echo "Docker is not running. Start Docker Desktop and retry."
  exit 1
fi

if docker ps --format '{{.Names}}' | grep -q '^finops-pg$'; then
  echo "WARNING: standalone container finops-pg is running on port 5432."
  echo "         Stop it to avoid conflicts: docker stop finops-pg"
fi

echo "==> starting core infrastructure (docker compose --profile core)"
docker compose --profile core up -d

echo "==> waiting for core health checks"
deadline=$((SECONDS + 180))
while (( SECONDS < deadline )); do
  if docker compose ps --status running | grep -q postgres \
    && docker compose ps --status running | grep -q redis \
    && docker compose ps --status running | grep -q kafka; then
    unhealthy="$(docker compose ps --format json | grep -c '"Health":"unhealthy"' || true)"
    if [[ "$unhealthy" == "0" ]]; then
      break
    fi
  fi
  sleep 5
done

bash "$ROOT/scripts/init-local.sh"

if command -v buf >/dev/null 2>&1; then
  echo "==> proto generate"
  pnpm proto:generate
else
  echo "==> skipping proto generate (buf not installed)"
fi

if command -v sqlc >/dev/null 2>&1; then
  echo "==> sqlc generate"
  pnpm sqlc:generate
else
  echo "==> skipping sqlc generate (sqlc not installed)"
fi

echo "Dev preflight complete."