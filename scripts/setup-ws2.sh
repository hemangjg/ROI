#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

if ! command -v buf >/dev/null 2>&1; then
  echo "Installing buf..."
  brew install bufbuild/buf/buf
fi

cd "$ROOT"
pnpm install

cd "$ROOT/contracts/proto"
buf lint
buf generate

cd "$ROOT/packages/proto"
go mod tidy
go test ./...
go build ./...

cd "$ROOT"
pnpm exec tsc -p packages/proto/tsconfig.json --noEmit
pnpm exec nx run workspace:proto-generate
pnpm exec nx run proto:build

bash "$ROOT/scripts/test-packages.sh"

echo "WS2 verification complete."
find "$ROOT/packages/proto/gen/go" "$ROOT/packages/proto/gen/ts" -type f 2>/dev/null | sort