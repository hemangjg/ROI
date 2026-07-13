#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"

if ! docker info >/dev/null 2>&1; then
  echo "Docker is not running. Start Docker Desktop and retry."
  exit 1
fi

echo "==> package tests"
bash "$ROOT/scripts/test-packages.sh"

echo "==> build all services"
bash "$ROOT/scripts/setup-ws5.sh"

echo "==> starting observability stack"
export OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
docker compose --profile obs up -d

echo "==> rebuilding services with OTLP enabled"
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

echo "==> generating traced requests"
for _ in {1..5}; do
  curl -fsS http://localhost:8080/ >/dev/null
  curl -fsS http://localhost:8081/ >/dev/null
done

echo "==> waiting for trace export"
sleep 10

echo "==> checking Jaeger for identity-service traces"
jaeger_deadline=$((SECONDS + 60))
jaeger_ready=false
while (( SECONDS < jaeger_deadline )); do
  if curl -fsS "http://localhost:16686/api/services" | grep -q identity-service; then
    jaeger_ready=true
    break
  fi
  sleep 5
done

if [[ "$jaeger_ready" != "true" ]]; then
  echo "ERROR: identity-service not found in Jaeger"
  curl -fsS "http://localhost:16686/api/services" || true
  exit 1
fi

echo "==> checking Prometheus scrape target"
if ! curl -fsS http://localhost:9090/-/healthy >/dev/null; then
  echo "ERROR: Prometheus is not healthy"
  exit 1
fi

echo "==> checking OTel collector Prometheus exporter"
collector_status="$(curl -fsS -o /dev/null -w "%{http_code}" http://localhost:8889/metrics)"
if [[ "$collector_status" != "200" ]]; then
  echo "ERROR: OTel collector metrics endpoint returned HTTP $collector_status"
  exit 1
fi

echo "WS7 verification complete."
echo "  Jaeger UI:    http://localhost:16686"
echo "  Prometheus:   http://localhost:9090"
echo "  Grafana:      http://localhost:3001 (admin/admin)"