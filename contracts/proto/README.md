# Protocol Buffers — Domain Contracts

Active protobuf contracts for the AI FinOps Ingestion MVP (Phase 1b).

## Layout

```
contracts/proto/
├── buf.yaml
├── buf.gen.yaml
├── bootstrap/
│   └── health/v1/health.proto       # toolchain verification
├── auth/v1/auth.proto               # identity-service (ValidateToken, ValidateApiKey)
├── ingestion/v1/event.proto         # ingestion-service + Kafka events
├── pricing/v1/pricing.proto         # pricing-service
├── analytics/v1/spend.proto         # analytics-service
├── management/v1/org.proto          # management-service
└── workflow/v1/workflow.proto       # workflow-service
```

Generated output:

```
packages/proto/gen/
├── go/                 # Go protobuf + gRPC stubs
└── ts/                 # TypeScript protobuf (protobuf-es)
```

## Commands

```bash
# Lint
pnpm proto:lint
# or: cd contracts/proto && buf lint

# Generate
pnpm proto:generate
# or: cd contracts/proto && buf generate

# Verify Phase 1b WS2
bash scripts/setup-p2-ws2.sh
# or: make p2-ws2

# Breaking change detection (against main)
cd contracts/proto && buf breaking --against '.git#branch=main'
```

## Nx targets

| Target | Project | Command |
|--------|---------|---------|
| `proto-lint` | workspace | `buf lint` |
| `proto-generate` | workspace | `buf lint` + `buf generate` |
| `build` | proto | compile Go module + TypeScript |

## Archived reference protos

Original Phase 0.5 drafts are preserved at
[`docs/reference/phase-0.5-proto/`](../../docs/reference/phase-0.5-proto/).
Active contracts in this directory may diverge (e.g. analytics RPC message names follow Buf STANDARD lint).

## Codegen drift

CI checks that generated output is current:

```bash
bash scripts/check-codegen-drift.sh
```

Requires at least one git commit; skipped on a fresh clone until the first commit.