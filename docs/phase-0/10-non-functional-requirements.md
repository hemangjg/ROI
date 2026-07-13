# 10 — Non-Functional Requirements

## Overview

Non-functional requirements define how the system performs, scales, and protects data — not what it does. These are contractual commitments for Enterprise tier and engineering targets for all tiers.

---

## Availability and Reliability

| ID | Requirement | Target | Tier |
|----|-------------|--------|------|
| NFR-AVAIL-001 | Platform uptime SLA | 99.9% (8.7h downtime/year) | Enterprise |
| NFR-AVAIL-002 | Platform uptime target | 99.5% | Growth/Starter |
| NFR-AVAIL-003 | Ingestion API availability | 99.95% | All |
| NFR-AVAIL-004 | Planned maintenance window | Sundays 02:00–06:00 UTC, notified 72h ahead | All |
| NFR-AVAIL-005 | Recovery Time Objective (RTO) | < 1 hour | Enterprise |
| NFR-AVAIL-006 | Recovery Point Objective (RPO) | < 5 minutes (event data) | All |
| NFR-AVAIL-007 | Graceful degradation: ingestion continues if dashboard is down | Required | All |
| NFR-AVAIL-008 | Graceful degradation: dashboard works if ingestion is delayed | Required (reads cached rollups) | All |

---

## Performance

| ID | Requirement | Target | Measurement |
|----|-------------|--------|-------------|
| NFR-PERF-001 | Ingestion API response time (p99) | < 100ms | APM |
| NFR-PERF-002 | Ingestion API response time (p50) | < 30ms | APM |
| NFR-PERF-003 | Event processing lag (ingest → queryable) | < 30s (p99) | Queue monitoring |
| NFR-PERF-004 | Dashboard page load (p95) | < 2s | RUM |
| NFR-PERF-005 | Spend summary API response (p95) | < 500ms | APM |
| NFR-PERF-006 | Time series query (90 days daily) | < 1s | APM |
| NFR-PERF-007 | Team spend table (50 teams) | < 1s | APM |
| NFR-PERF-008 | Batch ingest (500 events) | < 200ms (p99) | APM |

---

## Scalability

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-SCALE-001 | Events ingested per day (single org) | 10,000,000 |
| NFR-SCALE-002 | Events ingested per day (platform total) | 500,000,000 |
| NFR-SCALE-003 | Concurrent dashboard users (platform) | 5,000 |
| NFR-SCALE-004 | Organizations on platform | 10,000 |
| NFR-SCALE-005 | Teams per organization | 500 |
| NFR-SCALE-006 | API keys per organization | 100 |
| NFR-SCALE-007 | Raw event retention (hot storage) | 90 days |
| NFR-SCALE-008 | Rollup retention | 3 years |
| NFR-SCALE-009 | Horizontal scaling for API and workers | Stateless, K8s HPA |
| NFR-SCALE-010 | Database partitioning strategy | Monthly partitions on usage_events at 10M rows |

---

## Security

| ID | Requirement | Detail |
|----|-------------|--------|
| NFR-SEC-001 | Encryption in transit | TLS 1.2+ on all endpoints |
| NFR-SEC-002 | Encryption at rest | AES-256 (cloud-managed or application-level) |
| NFR-SEC-003 | API key storage | bcrypt hash; plaintext shown once at creation |
| NFR-SEC-004 | Password storage | bcrypt with cost factor 12 |
| NFR-SEC-005 | Tenant isolation | Row-level org_id filter on every query |
| NFR-SEC-006 | No cross-tenant data leakage | Integration tests validating isolation |
| NFR-SEC-007 | Secrets management | Environment variables via secrets manager (not in code) |
| NFR-SEC-008 | Input validation | All API inputs validated (Zod/Joi schemas) |
| NFR-SEC-009 | Rate limiting | Per API key and per IP |
| NFR-SEC-010 | CORS policy | Restricted to known dashboard domains |
| NFR-SEC-011 | Dependency scanning | Automated CVE scanning in CI |
| NFR-SEC-012 | Penetration testing | Annual third-party pentest (Enterprise launch) |
| NFR-SEC-013 | No PII in Phase 1 | Prompt content not stored; metadata tags only |
| NFR-SEC-014 | Audit log immutability | Append-only audit_logs table; no UPDATE/DELETE |

---

## Compliance

| ID | Requirement | Timeline |
|----|-------------|----------|
| NFR-COMP-001 | GDPR-ready data processing | Phase 1 |
| NFR-COMP-002 | Data Processing Agreement (DPA) template | Phase 2 |
| NFR-COMP-003 | SOC 2 Type II certification | Phase 4 (12 months post-launch) |
| NFR-COMP-004 | Right to erasure (GDPR Article 17) | Phase 4 |
| NFR-COMP-005 | Data residency (EU region option) | Phase 4 |
| NFR-COMP-006 | HIPAA BAA availability | Future (healthcare vertical) |

---

## Observability (Internal)

| ID | Requirement | Detail |
|----|-------------|--------|
| NFR-OBS-001 | Structured JSON logging | All services |
| NFR-OBS-002 | Request tracing (correlation ID) | End-to-end ingest → process → store |
| NFR-OBS-003 | Metrics export | Prometheus-compatible |
| NFR-OBS-004 | Error alerting | PagerDuty/Opsgenie for P1 incidents |
| NFR-OBS-005 | Ingest volume dashboards | Events/sec, queue depth, processing lag |
| NFR-OBS-006 | Cost calculation accuracy monitoring | Daily reconciliation vs. provider invoices |

---

## Deployability

| ID | Requirement | Detail |
|----|-------------|--------|
| NFR-DEPLOY-001 | Local development via Docker Compose | Single command: `docker compose up` |
| NFR-DEPLOY-002 | CI/CD pipeline | GitHub Actions: lint, test, build, deploy |
| NFR-DEPLOY-003 | Kubernetes-ready | Helm charts or K8s manifests |
| NFR-DEPLOY-004 | Blue-green or rolling deployments | Zero-downtime deploys |
| NFR-DEPLOY-005 | Database migrations | Automated via Prisma Migrate or similar |
| NFR-DEPLOY-006 | Environment separation | dev, staging, production |
| NFR-DEPLOY-007 | Infrastructure as Code | Terraform or Pulumi for cloud resources |

---

## Maintainability

| ID | Requirement | Detail |
|----|-------------|--------|
| NFR-MAINT-001 | TypeScript strict mode | All application code |
| NFR-MAINT-002 | Test coverage | > 80% on cost calculation and auth modules |
| NFR-MAINT-003 | API versioning | /v1/ prefix; breaking changes require /v2/ |
| NFR-MAINT-004 | OpenAPI specification | Auto-generated from code |
| NFR-MAINT-005 | Monorepo with shared types | packages/shared for schemas and types |
| NFR-MAINT-006 | Architecture Decision Records | docs/adr/ for irreversible decisions |
| NFR-MAINT-007 | Code review required | All changes via PR |

---

## Disaster Recovery

| ID | Requirement | Detail |
|----|-------------|--------|
| NFR-DR-001 | Database backups | Automated daily; 30-day retention |
| NFR-DR-002 | Point-in-time recovery | PostgreSQL WAL archiving |
| NFR-DR-003 | Multi-AZ deployment | Production database and API |
| NFR-DR-004 | Backup restoration tested | Quarterly DR drill |
| NFR-DR-005 | Event replay capability | Reprocess from queue dead-letter or raw archive |