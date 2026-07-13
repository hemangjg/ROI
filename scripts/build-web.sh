#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT/apps/web"

for attempt in 1 2 3 4 5; do
  if pnpm build; then
    exit 0
  fi

  echo "Web build attempt ${attempt} failed."
  if (( attempt < 5 )); then
    echo "Retrying after regenerating partial .next artifacts..."
    sleep 2
  fi
done

echo "ERROR: web build failed after 5 attempts"
exit 1