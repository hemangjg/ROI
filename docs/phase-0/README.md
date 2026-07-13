# Phase 0 — Product Definition

This directory contains the complete Phase 0 product definition for the AI FinOps platform. These documents define **what** to build, **for whom**, and **why** — before any technical architecture or code.

## Reading Order

| # | Document | Purpose |
|---|----------|---------|
| 01 | [Problem Statement](./01-problem-statement.md) | Why this product exists |
| 02 | [Target Customers](./02-target-customers.md) | Who we sell to and ICP |
| 03 | [Competitor Analysis](./03-competitor-analysis.md) | Competitive landscape and positioning |
| 04 | [Unique Selling Proposition](./04-unique-selling-proposition.md) | Why us, not alternatives |
| 05 | [Pricing Model](./05-pricing-model.md) | Tiers, unit economics, expansion |
| 06 | [Success Metrics](./06-success-metrics.md) | KPIs for product, business, and ops |
| 07 | [User Personas](./07-user-personas.md) | Alex, Morgan, Sam, Jordan, Riley |
| 08 | [Customer Journey](./08-customer-journey.md) | Awareness → Enterprise lifecycle |
| 09 | [Functional Requirements](./09-functional-requirements.md) | MoSCoW requirements by product area |
| 10 | [Non-Functional Requirements](./10-non-functional-requirements.md) | Performance, security, scale, compliance |
| 11 | [Risks and Assumptions](./11-risks-and-assumptions.md) | Risk register and validation plan |
| 12 | [Future Roadmap](./12-future-roadmap.md) | Phase 0.5 through Phase 4 plan |

## Phase 0 Exit Criteria

Before proceeding to Phase 0.5 (Technical Architecture):

- [ ] 5 customer discovery interviews completed
- [ ] 3+ prospects confirm pain and willingness to pay
- [ ] Founding team aligned on ICP, pricing, and Phase 1 scope
- [ ] Go/no-go decision documented

## What Comes Next

**Phase 0.5 — Technical Architecture** produces:
- Data model and database schema
- API contracts (OpenAPI)
- System design (ingestion pipeline, queue, workers)
- Security and tenancy design
- Monorepo structure
- Architecture Decision Records

**Phase 1 — Ingestion MVP** implements the data spine and first dashboard.