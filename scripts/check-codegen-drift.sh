#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

cd "$ROOT"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "Not a git repository; skipping codegen drift check."
  exit 0
fi

if ! git rev-parse --verify HEAD >/dev/null 2>&1; then
  echo "No git commits yet; skipping codegen drift check."
  exit 0
fi

if command -v sqlc >/dev/null 2>&1; then
  echo "==> sqlc generate"
  pnpm sqlc:generate
else
  echo "sqlc not installed; skipping sqlc drift check."
fi

if [[ -n "$(git status --porcelain db/atlas/sql/gen)" ]]; then
  echo "ERROR: SQLC generated code is out of date. Run pnpm sqlc:generate."
  git status --short db/atlas/sql/gen
  exit 1
fi

echo "Codegen drift check passed."