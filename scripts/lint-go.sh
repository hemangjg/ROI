#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"
go work sync

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