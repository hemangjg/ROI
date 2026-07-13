#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"

if ! docker info >/dev/null 2>&1; then
  echo "Docker is not running. Start Docker Desktop and retry."
  exit 1
fi

echo "==> generating local TLS certs"
bash "$ROOT/scripts/gen-local-certs.sh"

echo "==> starting backend services"
docker compose --profile services up -d --build

echo "==> waiting for envoy admin"
envoy_deadline=$((SECONDS + 120))
envoy_ready=false
while (( SECONDS < envoy_deadline )); do
  if curl -fsS http://127.0.0.1:9901/ready 2>/dev/null | grep -q LIVE; then
    envoy_ready=true
    break
  fi
  sleep 3
done

if [[ "$envoy_ready" != "true" ]]; then
  echo "ERROR: Envoy admin did not become ready"
  docker compose logs envoy --tail 50
  exit 1
fi

echo "==> gateway health"
curl -fsS http://localhost:8888/gateway/healthz | grep -q '"status":"ok"'

echo "==> HTTP route proxies"
for svc in identity ingestion management analytics pricing workflow; do
  body="$(curl -fsS "http://localhost:8888/_gateway/healthz/${svc}")"
  echo "    /_gateway/healthz/${svc} -> ${body}"
  echo "$body" | grep -q '"status":"ok"'
done

echo "==> OpenAPI path routing (404 expected until Phase 2 handlers exist)"
for path in /v1/events /v1/orgs /v1/orgs/demo/spend; do
  code="$(curl -sS -o /dev/null -w "%{http_code}" "http://localhost:8888${path}")"
  echo "    ${path} -> HTTP ${code}"
  if [[ "$code" != "404" && "$code" != "405" && "$code" != "200" ]]; then
    echo "ERROR: unexpected status ${code} for ${path}"
    exit 1
  fi
done

echo "==> HTTPS API listener"
curl -kfsS --resolve api.ai-finops.local:8443:127.0.0.1 \
  https://api.ai-finops.local:8443/gateway/healthz | grep -q '"status":"ok"'

echo "==> HTTPS app listener (web profile)"
docker compose --profile web up -d web
app_deadline=$((SECONDS + 60))
app_ready=false
while (( SECONDS < app_deadline )); do
  if curl -kfsS --resolve app.ai-finops.local:8444:127.0.0.1 \
    https://app.ai-finops.local:8444/ 2>/dev/null | grep -q "AI FinOps"; then
    app_ready=true
    break
  fi
  sleep 3
done

if [[ "$app_ready" != "true" ]]; then
  echo "ERROR: app.ai-finops.local did not route to web placeholder"
  exit 1
fi

echo "WS8 verification complete."
echo "  HTTP API:   http://localhost:8888/v1"
echo "  HTTPS API:  https://api.ai-finops.local:8443/v1  (add to /etc/hosts or use --resolve)"
echo "  HTTPS App:  https://app.ai-finops.local:8444      (web profile)"
echo "  Envoy admin: http://127.0.0.1:9901"