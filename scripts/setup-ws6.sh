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

echo "==> docker compose config"
docker compose config >/dev/null

echo "==> starting core infrastructure"
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

echo "==> building and starting application services"
docker compose --profile services up -d --build

echo "==> waiting for service health checks"
service_deadline=$((SECONDS + 300))
ready=false
while (( SECONDS < service_deadline )); do
  if curl -fsS http://localhost:8080/healthz >/dev/null \
    && curl -fsS http://localhost:8081/healthz >/dev/null \
    && curl -fsS http://localhost:8082/healthz >/dev/null \
    && curl -fsS http://localhost:8083/healthz >/dev/null \
    && curl -fsS http://localhost:8084/healthz >/dev/null \
    && curl -fsS http://localhost:8085/healthz >/dev/null; then
    ready=true
    break
  fi
  sleep 5
done

if [[ "$ready" != "true" ]]; then
  echo "ERROR: not all services became healthy in time"
  docker compose ps
  exit 1
fi

echo "==> health checks"
for port in 8080 8081 8082 8083 8084 8085; do
  echo "    :$port/healthz -> $(curl -fsS "http://localhost:${port}/healthz")"
  echo "    :$port/readyz  -> $(curl -fsS "http://localhost:${port}/readyz")"
done

echo "WS6 verification complete."