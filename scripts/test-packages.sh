#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"
go work sync

for mod in packages/*/; do
  if [[ -f "${mod}go.mod" ]]; then
    echo "==> testing ${mod}"
    (cd "$mod" && go test ./...)
  fi
done

if [[ -f db/atlas/sql/go.mod ]]; then
  echo "==> testing db/atlas/sql"
  (cd db/atlas/sql && go test ./... -count=1)
fi

echo "All package tests passed."