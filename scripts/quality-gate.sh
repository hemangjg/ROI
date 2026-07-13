#!/usr/bin/env bash
# Local quality gate for the concrete engineering base.
# Usage:
#   bash scripts/quality-gate.sh           # offline CI
#   WITH_E2E=1 bash scripts/quality-gate.sh  # + live stack e2e
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:$(go env GOPATH)/bin:${PATH:-}"
export NX_DAEMON=false

cd "$ROOT"

echo "==> quality-gate: format"
pnpm format:check

echo "==> quality-gate: ci-local"
bash scripts/ci-local.sh

if [[ "${WITH_E2E:-0}" == "1" ]]; then
  if ! docker info >/dev/null 2>&1; then
    echo "ERROR: WITH_E2E=1 requires Docker" >&2
    exit 1
  fi
  echo "==> quality-gate: e2e-smoke"
  if ! curl -fsS http://localhost:8888/gateway/healthz >/dev/null 2>&1; then
    echo "    starting services profile..."
    docker compose --profile services up -d
    bash scripts/init-local.sh
    # wait for envoy
    for _ in $(seq 1 60); do
      curl -fsS http://localhost:8888/gateway/healthz >/dev/null 2>&1 && break
      sleep 2
    done
  fi
  bash scripts/e2e-smoke.sh
fi

echo ""
echo "Quality gate PASSED"
echo "  offline CI: ok"
if [[ "${WITH_E2E:-0}" == "1" ]]; then
  echo "  e2e smoke: ok"
fi
