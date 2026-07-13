# 11 — Risks and Assumptions

## Risk Register

### Critical Risks (High Impact, Address Immediately)

#### R-001: Customers Will Not Install SDK
| Attribute | Detail |
|-----------|--------|
| **Description** | Target customers resist modifying application code to add cost tracking SDK |
| **Likelihood** | Medium |
| **Impact** | Critical — no SDK means no data means no product value |
| **Mitigation** | (1) Make SDK integration < 30 minutes with copy-paste docs. (2) Phase 2: reverse proxy mode requiring zero code changes. (3) Phase 2: provider billing API import as supplement. (4) Validate in 5 customer interviews before building. |
| **Owner** | Product + Engineering |
| **Status** | Open — validate in customer discovery |

#### R-002: Incumbents Add Cross-Provider AI FinOps
| Attribute | Detail |
|-----------|--------|
| **Description** | CloudZero, Datadog, or Vantage ship comprehensive AI cost management module |
| **Likelihood** | High |
| **Impact** | High — reduces differentiation, increases CAC |
| **Mitigation** | (1) Move fast on governance + optimization depth (not just dashboards). (2) Target mid-market where incumbents are overpriced. (3) Build FinOps-native workflows incumbents treat as bolt-ons. (4) Establish customer data moat (historical spend, trends). |
| **Owner** | Product |
| **Status** | Open — monitor competitor announcements monthly |

#### R-003: Pricing Accuracy Drift
| Attribute | Detail |
|-----------|--------|
| **Description** | Provider pricing changes frequently; our cost calculations diverge from actual invoices |
| **Likelihood** | Medium |
| **Impact** | High — destroys trust with FinOps persona (primary buyer) |
| **Mitigation** | (1) Versioned pricing catalog with effective dates. (2) Weekly automated pricing sync. (3) Admin alert when catalog is stale > 7 days. (4) Monthly reconciliation report comparing platform totals to provider invoices. (5) Allow customer pricing overrides for enterprise discounts. |
| **Owner** | Engineering |
| **Status** | Open — design in Phase 0.5 |

---

### High Risks

#### R-004: Enterprise Sales Cycle Exceeds 6 Months
| Attribute | Detail |
|-----------|--------|
| **Description** | Enterprise deals require security review, procurement, legal — delaying revenue |
| **Likelihood** | High |
| **Impact** | Medium — slows revenue but does not kill product |
| **Mitigation** | (1) Bottom-up PLG motion: engineer installs SDK → FinOps sees value → converts. (2) Target Growth segment first (2–6 week sales cycle). (3) Enterprise features (SSO, SCIM) only after 30+ paying customers. |
| **Owner** | Sales/Product |
| **Status** | Open |

#### R-005: Low Initial Event Volume Per Customer
| Attribute | Detail |
|-----------|--------|
| **Description** | Customers integrate SDK but send few events; dashboard looks empty |
| **Likelihood** | Medium |
| **Impact** | Medium — trial does not reach "aha moment" |
| **Mitigation** | (1) Onboarding checklist: integrate at least 2 services. (2) Day 3 automated check: if < 100 events, trigger support outreach. (3) Demo/synthetic data mode for sales demos. (4) Integration health dashboard showing coverage %. |
| **Owner** | Customer Success |
| **Status** | Open |

#### R-006: Data Volume Exceeds Architecture Capacity
| Attribute | Detail |
|-----------|--------|
| **Description** | High-volume customers (10M+ events/day) overwhelm ingestion or storage |
| **Likelihood** | Low (early), Medium (Year 2) |
| **Impact** | High — service degradation affects all customers |
| **Mitigation** | (1) Async queue architecture from day one. (2) Monthly table partitioning. (3) Rollup-based dashboards (never scan raw events). (4) Load testing before launch. (5) Per-org rate limiting. |
| **Owner** | Engineering |
| **Status** | Open — address in Phase 0.5 architecture |

---

### Medium Risks

#### R-007: Open Source Competitor Gains Traction
| Attribute | Detail |
|-----------|--------|
| **Description** | Langfuse, Helicone, or new OSS project adds FinOps features and undercuts pricing |
| **Likelihood** | Medium |
| **Impact** | Medium — compresses pricing, increases support burden |
| **Mitigation** | Compete on FinOps depth, enterprise features, and support — not on price alone. OSS lacks SSO, SCIM, SLA, and dedicated support. |
| **Owner** | Product |
| **Status** | Open |

#### R-008: AI Spend Growth Slows
| Attribute | Detail |
|-----------|--------|
| **Description** | AI API spend plateaus or declines due to model efficiency gains or shift to self-hosted |
| **Likelihood** | Low (2026–2027), Medium (2028+) |
| **Impact** | High — reduces TAM and urgency |
| **Mitigation** | (1) Expand to governance and productivity value props (not just cost). (2) Support self-hosted (Ollama) tracking. (3) Pivot to AI ROI/productivity platform if cost urgency declines. |
| **Owner** | Product |
| **Status** | Monitor |

#### R-009: Key Person Dependency
| Attribute | Detail |
|-----------|--------|
| **Description** | Founding team lacks FinOps domain expertise or enterprise sales experience |
| **Likelihood** | Medium |
| **Impact** | Medium — slows product-market fit |
| **Mitigation** | (1) Conduct 10+ customer discovery interviews in Phase 0. (2) Hire or advise with FinOps practitioner. (3) Design partner program for first 5 customers. |
| **Owner** | Leadership |
| **Status** | Open |

#### R-010: Regulatory Changes Restrict AI Usage Tracking
| Attribute | Detail |
|-----------|--------|
| **Description** | Privacy regulations limit ability to track AI usage metadata |
| **Likelihood** | Low |
| **Impact** | Medium — may require architecture changes |
| **Mitigation** | (1) No prompt content storage in Phase 0–1. (2) Metadata-only tracking. (3) GDPR-ready design from day one. (4) Legal review before EU launch. |
| **Owner** | Legal/Product |
| **Status** | Open |

---

## Assumptions

### Market Assumptions

| ID | Assumption | Validation method | Status |
|----|------------|-------------------|--------|
| A-001 | Companies with > $50K/yr AI spend will pay $500–$5K/mo for visibility | 5 pricing validation interviews | Unvalidated |
| A-002 | AI API spend continues growing > 50% YoY through 2027 | Industry reports, earnings calls | Likely true |
| A-003 | 30% of mid-market SaaS companies have material AI spend today | Market research, LinkedIn job postings | Unvalidated |
| A-004 | FinOps persona exists at target companies with budget authority | Customer interviews | Unvalidated |
| A-005 | Cross-provider view is top-3 pain point (not just nice-to-have) | Customer interviews | Unvalidated |

### Product Assumptions

| ID | Assumption | Validation method | Status |
|----|------------|-------------------|--------|
| A-006 | SDK integration is acceptable for engineering teams | Customer interviews + trial data | Unvalidated |
| A-007 | OpenAI + Anthropic cover 70%+ of enterprise LLM API spend | Provider market share data | Likely true |
| A-008 | Time to first insight < 24 hours is achievable | Phase 1 build + design partner test | Unvalidated |
| A-009 | Cost calculation accuracy > 99% is achievable with versioned pricing | Phase 1 reconciliation test | Unvalidated |
| A-010 | Governance features accelerate enterprise deals (CISO involvement) | Sales data from Phase 2+ | Unvalidated |

### Business Assumptions

| ID | Assumption | Validation method | Status |
|----|------------|-------------------|--------|
| A-011 | Bottom-up PLG motion works for this product category | Trial conversion data | Unvalidated |
| A-012 | 40% trial-to-paid conversion is achievable | Trial conversion data | Unvalidated |
| A-013 | Net revenue retention > 120% is achievable via tier upgrades | Expansion revenue data | Unvalidated |
| A-014 | CAC < $5K via content marketing + PLG (no enterprise sales team initially) | Marketing analytics | Unvalidated |
| A-015 | Gross margin > 80% at scale | Infrastructure cost modeling | Likely true |

---

## Risk Review Cadence

| Activity | Frequency | Participants |
|----------|-----------|-------------|
| Risk register review | Monthly | Founding team |
| Assumption validation check | Per customer interview | Product |
| Competitive landscape scan | Monthly | Product |
| Architecture risk review | Per phase transition | Engineering |
| Financial risk review | Quarterly | Leadership |