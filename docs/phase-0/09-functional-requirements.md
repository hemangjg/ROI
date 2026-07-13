# 09 — Functional Requirements

## Requirement Organization

Requirements are organized by product area and prioritized using MoSCoW:

- **Must** — MVP / Phase 1. Product cannot launch without these.
- **Should** — Phase 2. Expected within 3 months of launch.
- **Could** — Phase 3+. Differentiating features.
- **Won't** — Explicitly out of scope for initial releases.

Each requirement has a unique ID for traceability: `FR-{area}-{number}`.

---

## Area 1: Organization and Identity

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-ORG-001 | Create organization with name and slug | Must | 1 |
| FR-ORG-002 | Invite users to organization via email | Must | 1 |
| FR-ORG-003 | Create teams within organization | Must | 1 |
| FR-ORG-004 | Assign users to teams with roles | Must | 1 |
| FR-ORG-005 | Soft-delete teams (preserve historical data) | Must | 1 |
| FR-ORG-006 | Organization settings page (name, timezone) | Must | 1 |
| FR-ORG-007 | SSO via SAML/OIDC | Should | 4 |
| FR-ORG-008 | SCIM user provisioning | Should | 4 |
| FR-ORG-009 | Multi-org parent/child hierarchy | Could | 4 |
| FR-ORG-010 | White-label branding (logo, colors) | Could | 4 |

---

## Area 2: Authentication and Authorization

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-AUTH-001 | Email/password registration and login | Must | 1 |
| FR-AUTH-002 | API key generation for ingestion | Must | 1 |
| FR-AUTH-003 | API key revocation | Must | 1 |
| FR-AUTH-004 | Role-based access control (admin, viewer) | Must | 1 |
| FR-AUTH-005 | API key scoped to organization | Must | 1 |
| FR-AUTH-006 | JWT session tokens for dashboard | Must | 1 |
| FR-AUTH-007 | API key shown only once at creation | Must | 1 |
| FR-AUTH-008 | Granular permissions (team-scoped viewer) | Should | 2 |
| FR-AUTH-009 | SSO with Okta, Azure AD, Google Workspace | Should | 4 |
| FR-AUTH-010 | MFA for dashboard users | Should | 2 |

---

## Area 3: Usage Event Ingestion

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-ING-001 | Accept single usage event via REST API | Must | 1 |
| FR-ING-002 | Accept batch of up to 500 events via REST API | Must | 1 |
| FR-ING-003 | Idempotency key deduplication | Must | 1 |
| FR-ING-004 | Async processing (202 Accepted response) | Must | 1 |
| FR-ING-005 | Event schema: provider, model, input_tokens, output_tokens, occurred_at, metadata | Must | 1 |
| FR-ING-006 | Metadata tags: team_id, user_id, project, feature, customer | Must | 1 |
| FR-ING-007 | Node.js SDK wrapping LLM client calls | Must | 1 |
| FR-ING-008 | Python SDK wrapping LLM client calls | Should | 2 |
| FR-ING-009 | Cache token fields (cache_read, cache_write) | Must | 1 |
| FR-ING-010 | Per-event cost calculation on ingest | Must | 1 |
| FR-ING-011 | Rate limiting per API key (1000 events/min) | Must | 1 |
| FR-ING-012 | Event status lookup by ID | Should | 1 |
| FR-ING-013 | Webhook notifications on processing failure | Could | 2 |
| FR-ING-014 | Reverse proxy / API gateway ingestion mode | Could | 3 |
| FR-ING-015 | Provider-native billing API import (supplement SDK) | Could | 2 |

---

## Area 4: Cost Calculation and Pricing

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-COST-001 | Pricing catalog for OpenAI models | Must | 1 |
| FR-COST-002 | Pricing catalog for Anthropic models | Must | 1 |
| FR-COST-003 | Versioned pricing with effective dates | Must | 1 |
| FR-COST-004 | Calculate cost from token counts + model pricing | Must | 1 |
| FR-COST-005 | Handle unknown models (flag as unpriced) | Must | 1 |
| FR-COST-006 | Support input, output, cache_read, cache_write token pricing | Must | 1 |
| FR-COST-007 | Admin alert when pricing catalog is stale (> 7 days) | Should | 1 |
| FR-COST-008 | Custom pricing override per org (enterprise discounts) | Should | 2 |
| FR-COST-009 | Pricing for Google Gemini, Grok, DeepSeek, Mistral | Should | 2 |
| FR-COST-010 | Pricing for Azure OpenAI, AWS Bedrock, Vertex AI | Could | 2 |
| FR-COST-011 | Zero-cost tracking for self-hosted (Ollama) | Could | 2 |
| FR-COST-012 | Batch pricing discounts applied automatically | Could | 3 |

---

## Area 5: Dashboards and Reporting

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-DASH-001 | Organization spend summary (today, MTD, vs last period) | Must | 1 |
| FR-DASH-002 | Spend time series chart (daily granularity) | Must | 1 |
| FR-DASH-003 | Spend breakdown by provider | Must | 1 |
| FR-DASH-004 | Spend breakdown by team (sortable table) | Must | 1 |
| FR-DASH-005 | Date range filter (7d, 30d, 90d, custom) | Must | 1 |
| FR-DASH-006 | Team filter (drill into single team) | Must | 1 |
| FR-DASH-007 | Model-level spend breakdown | Should | 2 |
| FR-DASH-008 | User-level spend breakdown | Should | 2 |
| FR-DASH-009 | Project/feature spend breakdown | Should | 2 |
| FR-DASH-010 | Export to CSV | Should | 2 |
| FR-DASH-011 | Executive summary PDF report | Could | 3 |
| FR-DASH-012 | Custom dashboard widgets | Could | 3 |
| FR-DASH-013 | Scheduled email reports | Could | 2 |
| FR-DASH-014 | Cost per customer (for SaaS companies) | Could | 3 |

---

## Area 6: Budgeting and Alerts

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-BUDG-001 | Set monthly budget per team | Should | 2 |
| FR-BUDG-002 | Soft limit alert at configurable threshold (e.g., 80%) | Should | 2 |
| FR-BUDG-003 | Hard limit blocking (reject events when budget exceeded) | Could | 3 |
| FR-BUDG-004 | Email and in-app notifications for alerts | Should | 2 |
| FR-BUDG-005 | Slack webhook integration for alerts | Could | 2 |
| FR-BUDG-006 | Spend forecasting (linear projection) | Should | 2 |
| FR-BUDG-007 | Budget vs. actual comparison view | Should | 2 |
| FR-BUDG-008 | Org-level budget (in addition to team) | Could | 2 |
| FR-BUDG-009 | Budget approval workflow | Could | 4 |

---

## Area 7: Governance

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-GOV-001 | Audit log for all admin actions | Must | 1 |
| FR-GOV-002 | Approved model allowlist per org | Could | 2 |
| FR-GOV-003 | Restricted model blocklist per org | Could | 2 |
| FR-GOV-004 | Team-level model policies | Could | 2 |
| FR-GOV-005 | Sensitive data detection in metadata | Could | 3 |
| FR-GOV-006 | Policy violation alerts | Could | 2 |
| FR-GOV-007 | Compliance report export (SOC 2, ISO 27001) | Could | 3 |
| FR-GOV-008 | Data retention policy configuration | Could | 3 |
| FR-GOV-009 | GDPR data deletion (right to erasure) | Should | 4 |

---

## Area 8: Optimization

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-OPT-001 | Detect expensive model usage where cheaper alternative exists | Could | 3 |
| FR-OPT-002 | Detect duplicate/similar prompts | Could | 3 |
| FR-OPT-003 | Recommend model downgrade (GPT-4 → GPT-4o-mini) | Could | 3 |
| FR-OPT-004 | Identify cache opportunities (repeated context) | Could | 3 |
| FR-OPT-005 | Batch processing recommendations | Could | 3 |
| FR-OPT-006 | Prompt optimization suggestions | Could | 3 |
| FR-OPT-007 | Monthly savings report | Could | 3 |
| FR-OPT-008 | ROI estimation (time saved vs. cost) | Could | 3 |

---

## Area 9: Provider Support

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-PROV-001 | OpenAI provider plugin | Must | 1 |
| FR-PROV-002 | Anthropic provider plugin | Must | 1 |
| FR-PROV-003 | Provider plugin interface (extensible) | Must | 1 |
| FR-PROV-004 | Google Gemini provider plugin | Should | 2 |
| FR-PROV-005 | Grok (xAI) provider plugin | Should | 2 |
| FR-PROV-006 | DeepSeek provider plugin | Should | 2 |
| FR-PROV-007 | Mistral provider plugin | Should | 2 |
| FR-PROV-008 | Azure OpenAI provider plugin | Could | 2 |
| FR-PROV-009 | AWS Bedrock provider plugin | Could | 2 |
| FR-PROV-010 | Vertex AI provider plugin | Could | 2 |
| FR-PROV-011 | Ollama (self-hosted) provider plugin | Could | 2 |
| FR-PROV-012 | Custom provider plugin SDK for customers | Could | 4 |

---

## Area 10: Administration

| ID | Requirement | Priority | Phase |
|----|-------------|----------|-------|
| FR-ADMIN-001 | API key management UI (create, list, revoke) | Must | 1 |
| FR-ADMIN-002 | User management UI (invite, remove, change role) | Must | 1 |
| FR-ADMIN-003 | Team management UI (create, edit, delete) | Must | 1 |
| FR-ADMIN-004 | Organization settings | Must | 1 |
| FR-ADMIN-005 | Audit log viewer | Must | 1 |
| FR-ADMIN-006 | Usage quota display (events used vs. limit) | Should | 1 |
| FR-ADMIN-007 | Billing and subscription management | Could | 4 |
| FR-ADMIN-008 | Data export (full org data) | Could | 3 |

---

## Phase 1 MVP Requirement Summary

Phase 1 must ship these 35 requirements:

```
FR-ORG-001, FR-ORG-002, FR-ORG-003, FR-ORG-004, FR-ORG-005, FR-ORG-006
FR-AUTH-001 through FR-AUTH-007
FR-ING-001 through FR-ING-007, FR-ING-009, FR-ING-010, FR-ING-011
FR-COST-001 through FR-COST-006
FR-DASH-001 through FR-DASH-006
FR-GOV-001
FR-PROV-001 through FR-PROV-003
FR-ADMIN-001 through FR-ADMIN-005
```