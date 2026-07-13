# ADR-003: Tenant Isolation with OpenFGA

## Status

Accepted

## Context

Enterprise customers require strict multi-tenant isolation and fine-grained authorization (org admin, team viewer, API key manager). Hardcoded RBAC middleware will not scale to team-scoped permissions and ABAC rules in Phase 2.

## Decision

1. **Row-level isolation:** every PostgreSQL and ClickHouse query includes `org_id` filter
2. **OpenFGA** for authorization decisions on Management and Analytics services
3. **Relationship tuples** model org/team/user permissions
4. **Envoy ext_authz** delegates JWT validation to identity-service
5. **Ingest path** validates API keys in Ingestion Service (not OpenFGA) for throughput

## Consequences

**Positive:**
- Fine-grained permissions without code changes (tuple updates)
- ABAC extension path for team-scoped spend visibility
- Industry-standard ReBAC model

**Negative:**
- Additional infrastructure (OpenFGA server)
- Authorization check latency (~5ms per request)
- Tuple management complexity on org/team changes

## Alternatives Considered

| Alternative | Rejected because |
|-------------|-----------------|
| Hardcoded RBAC | Cannot express team-scoped viewer without code changes |
| OPA only | Better for policy rules; OpenFGA better for relationship-based access |
| Schema-per-tenant | Operational nightmare at 10,000 orgs |
| Database-per-tenant | Cost prohibitive; connection pool exhaustion |