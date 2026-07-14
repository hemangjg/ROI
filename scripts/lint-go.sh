#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"
go work sync

# Generated protos are gitignored; ensure they exist before analysis.
if [[ ! -d packages/proto/gen/go ]]; then
  if command -v buf >/dev/null 2>&1 && command -v pnpm >/dev/null 2>&1; then
    echo "==> packages/proto/gen missing — running pnpm proto:generate"
    pnpm proto:generate
  else
    echo "ERROR: packages/proto/gen is missing and buf/pnpm unavailable to generate it." >&2
    echo "Install buf + pnpm, then run: pnpm proto:generate" >&2
    exit 1
  fi
fi

lint_module() {
  local mod="$1"
  if [[ ! -f "${mod}/go.mod" ]]; then
    return 0
  fi

  if command -v golangci-lint >/dev/null 2>&1; then
    echo "    golangci-lint ${mod}"
    (cd "$mod" && golangci-lint run ./...)
  else
    echo "    go vet ${mod}"
    (cd "$mod" && go vet ./...)
  fi
}

if command -v golangci-lint >/dev/null 2>&1; then
  echo "==> golangci-lint (per module)"
else
  echo "==> golangci-lint not installed; falling back to go vet"
fi

for mod in packages/*/ services/*/ db/atlas/sql; do
  lint_module "$mod"
done

echo "Go lint complete."
