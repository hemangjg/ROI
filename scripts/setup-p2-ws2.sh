#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"
export NX_DAEMON=false

cd "$ROOT"

echo "==> verify proto artifacts"
for file in \
  contracts/proto/auth/v1/auth.proto \
  contracts/proto/ingestion/v1/event.proto \
  contracts/proto/pricing/v1/pricing.proto \
  contracts/proto/analytics/v1/spend.proto \
  contracts/proto/management/v1/org.proto \
  contracts/proto/workflow/v1/workflow.proto; do
  test -f "$file"
  echo "    ok: $file"
done

if ! command -v buf >/dev/null 2>&1; then
  echo "ERROR: buf not installed (brew install bufbuild/buf/buf)"
  exit 1
fi

cd "$ROOT/contracts/proto"
echo "==> buf lint"
buf lint

echo "==> buf generate"
buf generate

echo "==> verify generated output"
for file in \
  packages/proto/gen/go/auth/v1/auth.pb.go \
  packages/proto/gen/go/auth/v1/auth_grpc.pb.go \
  packages/proto/gen/go/ingestion/v1/event.pb.go \
  packages/proto/gen/go/ingestion/v1/event_grpc.pb.go \
  packages/proto/gen/go/pricing/v1/pricing.pb.go \
  packages/proto/gen/go/pricing/v1/pricing_grpc.pb.go \
  packages/proto/gen/go/analytics/v1/spend.pb.go \
  packages/proto/gen/go/analytics/v1/spend_grpc.pb.go \
  packages/proto/gen/go/management/v1/org.pb.go \
  packages/proto/gen/go/management/v1/org_grpc.pb.go \
  packages/proto/gen/go/workflow/v1/workflow.pb.go \
  packages/proto/gen/go/workflow/v1/workflow_grpc.pb.go \
  packages/proto/gen/ts/auth/v1/auth_pb.ts \
  packages/proto/gen/ts/ingestion/v1/event_pb.ts \
  packages/proto/gen/ts/pricing/v1/pricing_pb.ts \
  packages/proto/gen/ts/analytics/v1/spend_pb.ts \
  packages/proto/gen/ts/management/v1/org_pb.ts \
  packages/proto/gen/ts/workflow/v1/workflow_pb.ts; do
  test -f "$ROOT/$file"
  echo "    ok: $file"
done

cd "$ROOT/packages/proto"
echo "==> go test"
go mod tidy
go test ./... -count=1
go build ./...

echo "==> typescript compile"
cd "$ROOT"
pnpm exec tsc -p packages/proto/tsconfig.json --noEmit
pnpm exec nx run proto:build

echo ""
echo "P2-WS2 exit criteria:"
echo "  [x] auth/v1/auth.proto (AuthService)"
echo "  [x] ingestion/v1/event.proto (IngestionService + UsageEventCreated)"
echo "  [x] pricing/v1/pricing.proto (PricingService)"
echo "  [x] analytics/v1/spend.proto (AnalyticsService)"
echo "  [x] management/v1/org.proto (ManagementService)"
echo "  [x] workflow/v1/workflow.proto (WorkflowService)"
echo "  [x] buf lint + generate (Go + TypeScript)"
echo ""
echo "P2-WS2 verification complete."
echo "  Next: P2-WS3 identity-service auth APIs"