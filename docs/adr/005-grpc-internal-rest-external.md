# ADR-005: gRPC Internal, REST External

## Status

Accepted

## Context

SDKs and the Next.js dashboard need simple HTTP integration. Internal service communication (Workflow → Pricing, Envoy → Auth) benefits from typed contracts and lower latency.

## Decision

1. **External APIs:** REST with OpenAPI 3.1 spec (`api/openapi/v1.yaml`)
2. **Internal APIs:** gRPC with Protocol Buffers (`api/proto/`, managed by Buf)
3. **Pricing and Workflow services:** gRPC only, no public HTTP routes
4. **Envoy ext_authz:** gRPC to identity-service

SDKs and dashboard never call gRPC directly.

## Consequences

**Positive:**
- SDK integration is standard HTTP (any language)
- Internal type safety via generated protobuf stubs
- Buf enforces breaking change detection in CI

**Negative:**
- Two API specs to maintain (mitigated: shared event schemas)
- gRPC debugging requires grpcurl or service mesh tooling

## Alternatives Considered

| Alternative | Rejected because |
|-------------|-----------------|
| REST everywhere | Slower internal communication; no strong typing |
| gRPC everywhere | SDK integration complexity; browser gRPC-web overhead |
| GraphQL | Over-engineered for Phase 1 read patterns |