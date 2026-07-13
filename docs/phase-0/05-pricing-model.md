# 05 — Pricing Model

## Pricing Philosophy

Price on value delivered, not on AI spend percentage. Customers resist sharing spend data for revenue-share models. Flat subscription + usage overage is predictable for both sides and aligns with how FinOps tools (CloudZero, Vantage) price.

**Core principle:** A customer saving $50K/year in AI waste should happily pay $24K/year for the platform.

---

## Tier Structure

### Starter — $499/month

| Included | Limit |
|----------|-------|
| Organizations | 1 |
| Teams | 5 |
| Users | 10 |
| Events per month | 1,000,000 |
| Providers | 2 (OpenAI + Anthropic) |
| Data retention | 90 days |
| Dashboards | Spend overview, by provider, by team |
| API keys | 3 |
| Support | Email, 48h response |

**Target:** SMB pilots, agencies, teams testing AI cost visibility.

---

### Growth — $1,999/month

| Included | Limit |
|----------|-------|
| Organizations | 1 |
| Teams | Unlimited |
| Users | 50 |
| Events per month | 10,000,000 |
| Providers | 5 |
| Data retention | 1 year |
| Dashboards | All Starter + forecasting, trends |
| Features | Team budgets, soft/hard alerts |
| API keys | 10 |
| Support | Email + chat, 24h response |

**Target:** Series B–D SaaS companies with material AI spend and multiple teams.

---

### Enterprise — Custom ($5,000–$25,000/month)

| Included | Limit |
|----------|-------|
| Organizations | Multiple (parent/child) |
| Teams | Unlimited |
| Users | Unlimited |
| Events per month | Unlimited |
| Providers | All supported + custom plugins |
| Data retention | Custom (up to 3 years) |
| Features | All Growth + governance, SSO, SCIM, audit, white-label |
| API keys | Unlimited |
| Support | Dedicated CSM, 4h SLA, Slack channel |
| SLA | 99.9% uptime |

**Target:** Enterprise 2,000+ employees, regulated industries, multi-division chargeback.

---

## Usage Overage

| Tier | Overage rate |
|------|-------------|
| Starter | $0.50 per 100K events beyond 1M |
| Growth | $0.40 per 100K events beyond 10M |
| Enterprise | Negotiated (typically $0.20–$0.30 per 100K) |

Overage is billed monthly in arrears. Customers receive alerts at 80% and 100% of event quota.

---

## Trial and Onboarding

| Element | Detail |
|---------|--------|
| **Free trial** | 14 days, Growth tier features |
| **Credit card** | Required at signup (reduces tire-kickers) |
| **Trial data** | Retained if customer converts; deleted after 30 days if not |
| **Onboarding** | Self-serve SDK docs + optional 30-min onboarding call (Growth+) |
| **Time to value** | First dashboard within 24 hours of SDK integration |

---

## Pricing Alternatives Considered and Rejected

| Model | Why rejected |
|-------|-------------|
| **% of AI spend tracked** | Requires trust; customers resist sharing total spend; hard to verify |
| **% of savings** | Savings attribution is subjective; creates disputes; sales friction |
| **Per-seat only** | Does not scale with usage; heavy users subsidized by light users |
| **Open source + paid cloud** | Slows enterprise sales; support burden on free users |
| **Freemium forever** | Attracts low-spend users who never convert; infrastructure cost |

---

## Unit Economics Target

| Metric | Target |
|--------|--------|
| Average ACV (Year 1) | $18,000 |
| Gross margin | > 80% |
| CAC payback | < 12 months |
| LTV:CAC ratio | > 3:1 |
| Net revenue retention | > 120% |
| Logo churn (annual) | < 10% |

### Cost to serve (per customer, monthly)

| Component | Starter | Growth | Enterprise |
|-----------|---------|--------|------------|
| Infrastructure | ~$15 | ~$80 | ~$500 |
| Support | ~$20 | ~$50 | ~$2,000 |
| **Total COGS** | **~$35** | **~$130** | **~$2,500** |
| **Gross margin** | **~93%** | **~93%** | **~90%** |

---

## Expansion Revenue Paths

| Trigger | Upsell |
|---------|--------|
| Team count exceeds tier limit | Starter → Growth |
| Event volume exceeds quota | Overage billing (automatic) |
| Requests SSO / governance | Growth → Enterprise |
| Adds 3+ more providers | Growth → Enterprise |
| Multi-division chargeback | Enterprise multi-org |
| Wants optimization features | Enterprise add-on (Phase 3) |

---

## Competitive Pricing Position

| Competitor | Comparable tier price | Our position |
|------------|----------------------|--------------|
| Helicone Pro | ~$200/mo | We are more expensive but serve FinOps, not just devs |
| CloudZero | $50K+/yr | We are 4–10x cheaper, AI-focused |
| Datadog LLM add-on | $10K+/yr add-on | We are standalone, no infra monitoring required |
| LangSmith | Usage-based, ~$500/mo at scale | Comparable price, different buyer |

**Positioning:** Premium to dev tools, discount to enterprise FinOps platforms.

---

## Pricing Experiments (Post-Launch)

- [ ] Annual prepay discount (15% off) — test in Month 3
- [ ] Usage-based tier between Starter and Growth ($999/mo, 5M events) — test if Starter → Growth gap is too wide
- [ ] "Savings guarantee" pilot for Enterprise — refund if < 10% waste identified in 90 days
- [ ] Partner/reseller pricing for cloud consultancies