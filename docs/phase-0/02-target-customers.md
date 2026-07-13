# 02 — Target Customers

## Market Segmentation

### Primary Segment: Growth-Stage SaaS (Series B–D)

| Attribute | Profile |
|-----------|---------|
| **Company size** | 200–2,000 employees |
| **AI spend** | $100K–$2M per year |
| **AI maturity** | Production LLM features + internal tooling + employee AI tools |
| **Providers** | 2–5 (typically OpenAI, Anthropic, plus one cloud-hosted) |
| **Buyer** | VP Engineering or CTO |
| **Champion** | Platform Engineering lead or FinOps Manager |
| **Budget authority** | Engineering budget ($500K–$5M tooling spend) |
| **Sales cycle** | 2–6 weeks (bottom-up trial → team expansion) |

**Why they buy:** AI spend just became material enough to require visibility. They have no internal tooling to build this. Engineering leadership needs to justify spend to the board.

**Example companies:** Vertical SaaS platforms, fintech infrastructure, developer tools, HR tech, legal tech — any company shipping AI-powered features to customers.

---

### Secondary Segment: Enterprise (2,000+ employees)

| Attribute | Profile |
|-----------|---------|
| **Company size** | 2,000–50,000 employees |
| **AI spend** | $2M–$20M per year |
| **AI maturity** | Enterprise-wide AI adoption with governance requirements |
| **Providers** | 5–10+ including Azure OpenAI, AWS Bedrock, Vertex AI |
| **Buyer** | FinOps Director + CISO (joint decision) |
| **Champion** | AI Platform team or Cloud Center of Excellence |
| **Budget authority** | Central IT / FinOps budget |
| **Sales cycle** | 3–9 months (security review, procurement, pilot) |

**Why they buy:** Compliance, audit, and multi-department chargeback requirements. Cannot rely on per-provider dashboards. Need SSO, SCIM, and audit trails.

**Example companies:** Financial services, healthcare, insurance, large retail, government contractors.

---

### Tertiary Segment: AI-Heavy Agencies and Consultancies

| Attribute | Profile |
|-----------|---------|
| **Company size** | 50–500 employees |
| **AI spend** | $50K–$500K per year (often passed through to clients) |
| **AI maturity** | AI is the product — every project uses LLMs |
| **Providers** | 2–4, rotated per client requirements |
| **Buyer** | COO or Head of Finance |
| **Champion** | Technical Director |
| **Budget authority** | Operations budget |
| **Sales cycle** | 1–4 weeks |

**Why they buy:** Need to bill AI costs accurately to clients. Margin depends on understanding per-project AI spend.

---

## Ideal First Customer Profile (ICP)

The company we should land first:

```
Company:     Series C SaaS, ~800 employees
AI spend:    $400K/year, growing 200% YoY
Providers:   OpenAI (production), Anthropic (internal), Copilot (dev)
Pain:        Finance flagged a 3x cost increase with zero attribution
Trigger:     Board asked CTO for AI ROI report — CTO has 2 weeks
Champion:    Senior Platform Engineer who owns internal tooling
Buyer:       VP Engineering with FinOps dotted-line responsibility
Integration: Can install an npm SDK in their Node.js monorepo within a day
```

---

## Customer Qualification Criteria

### Must-have (deal qualifies)

- [ ] Annual LLM API spend > $50K (or projected to reach within 6 months)
- [ ] Using 2+ LLM providers OR planning to within 12 months
- [ ] Engineering team capable of SDK integration (or willing to use proxy mode later)
- [ ] A person exists who owns cloud/AI cost visibility (even informally)

### Nice-to-have (accelerates deal)

- [ ] Dedicated FinOps function
- [ ] Board or investor pressure on AI spend
- [ ] Existing cloud cost management tool (CloudZero, Kubecost, Vantage)
- [ ] SOC 2 or ISO 27001 certification (governance urgency)
- [ ] Node.js or Python primary stack (SDK compatibility)

### Disqualifiers

- AI spend < $10K/year (not enough pain to justify subscription)
- Single provider with no plans to diversify (provider dashboard suffices)
- No engineering resources for integration (cannot deliver value)
- Expecting fully passive monitoring with zero integration (not possible in Phase 1)

---

## Total Addressable Market (TAM) Estimate

| Segment | Companies (global) | Avg ACV | Segment TAM |
|---------|-------------------|---------|-------------|
| Growth SaaS (B–D) | ~15,000 | $24K/yr | $360M |
| Enterprise (2K+) | ~8,000 | $120K/yr | $960M |
| Agencies | ~5,000 | $12K/yr | $60M |
| **Total** | **~28,000** | — | **~$1.4B** |

Assumptions: 30% of companies in these segments have material AI spend today; penetration grows as AI adoption broadens. TAM expands to $5B+ by 2028 as AI spend becomes universal.

---

## Geographic Focus (Phase 0–1)

| Priority | Region | Rationale |
|----------|--------|-----------|
| 1 | North America | Highest AI adoption, strongest FinOps culture |
| 2 | UK / Western Europe | GDPR governance urgency, growing AI spend |
| 3 | ANZ | Early adopters, English-language sales |
| Defer | APAC, LATAM | Phase 3+ expansion |

---

## Anti-Personas (Who We Do Not Target Initially)

| Anti-persona | Why not |
|--------------|---------|
| Solo developers / indie hackers | Spend too low; use free tiers |
| Companies with < $10K AI spend | Pain insufficient for paid product |
| Pure research labs | Different buying process; academic pricing expectations |
| Companies building their own AI platform | Would build in-house; long sales cycle, low conversion |
| Non-technical businesses using only ChatGPT Enterprise | No API integration path; different product needed |