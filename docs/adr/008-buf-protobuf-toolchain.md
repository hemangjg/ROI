# ADR-008: Buf for Protocol Buffer Toolchain

## Status

Accepted

## Context

Six gRPC services share internal contracts. Ad-hoc `protoc` usage leads to inconsistent generation, no breaking change detection, and version drift across services.

## Decision

Use **Buf** as the protobuf toolchain:
- Module root: `api/proto/`
- `buf lint` enforces STANDARD lint rules in CI
- `buf breaking --against main` prevents breaking changes
- `buf generate` produces Go stubs into `gen/go/`
- Remote plugins: `protocolbuffers/go`, `grpc/go`

## Consequences

**Positive:**
- Consistent codegen across all services
- Breaking change detection in PR checks
- Remote plugins eliminate local protoc installation
- Managed mode sets consistent `go_package` prefixes

**Negative:**
- Buf CLI dependency for all developers
- Generated code in `gen/go/` must be gitignored or committed (decision: gitignored, generated in CI)

## Alternatives Considered

| Alternative | Rejected because |
|-------------|-----------------|
| Raw protoc | No lint/breaking checks; local plugin version drift |
| Connect RPC | Less ecosystem maturity; team chose standard gRPC |
| JSON over HTTP internal | No type safety; slower than protobuf |