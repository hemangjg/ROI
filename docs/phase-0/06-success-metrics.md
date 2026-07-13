# 06 — Success Metrics (KPIs)

## Metrics Framework

Metrics split into three layers: customer outcomes (does the product work?), business outcomes (does the company grow?), and operational health (can we deliver reliably?).

---

## Layer 1: Customer Outcome KPIs

These measure whether customers get value from the platform.

| KPI | Definition | Target | Measurement |
|-----|------------|--------|-------------|
| **Time to first insight** | Hours from SDK install to first dashboard with real data | < 24 hours | Onboarding funnel tracking |
| **Cost visibility coverage** | % of customer's total AI spend captured by platform | > 90% within 30 days | Customer-reported total vs. platform total |
| **Attribution completeness** | % of events with team/project metadata tags | > 85% | Event metadata analysis |
| **Waste identified** | Monthly dollar value of flagged optimization opportunities | 15–25% of tracked spend | Optimization engine output |
| **Budget compliance rate** | % of teams staying within soft budget limits | > 80% | Budget vs. actual comparison |
| **Alert action rate** | % of budget alerts that result in customer action within 7 days | > 50% | Alert → dashboard visit correlation |
| **Dashboard weekly active users** | Unique users viewing dashboards per week per org | > 3 per org | Product analytics |
| **NPS** | Net Promoter Score | > 40 | Quarterly survey |

---

## Layer 2: Business KPIs

These measure company growth and sustainability.

| KPI | Definition | Target (Year 1) | Target (Year 2) |
|-----|------------|-----------------|-----------------|
| **MRR** | Monthly recurring revenue | $50K | $300K |
| **ARR** | Annual recurring revenue | $600K | $3.6M |
| **Customers** | Paying organizations | 30 | 150 |
| **ACV** | Average contract value | $18K | $24K |
| **Trial-to-paid conversion** | % of trials converting to paid | > 40% | > 50% |
| **Logo churn (monthly)** | % of customers canceling per month | < 2% | < 1.5% |
| **Net revenue retention** | Revenue from existing customers (expansion - churn) | > 110% | > 120% |
| **CAC** | Customer acquisition cost | < $5K | < $8K |
| **CAC payback** | Months to recover acquisition cost | < 12 | < 10 |
| **LTV:CAC** | Lifetime value to acquisition cost ratio | > 3:1 | > 4:1 |
| **Gross margin** | Revenue minus COGS | > 80% | > 85% |

---

## Layer 3: Operational KPIs

These measure platform reliability and engineering health.

| KPI | Definition | Target |
|-----|------------|--------|
| **Ingest uptime** | % of time ingestion API accepts events | 99.95% |
| **Ingest latency (p99)** | 99th percentile API response time | < 100ms |
| **Event processing lag** | Time from ingest to cost-calculated and queryable | < 30 seconds (p99) |
| **Dashboard load time (p95)** | 95th percentile page load for spend queries | < 2 seconds |
| **Cost calculation accuracy** | % of events with correct cost vs. provider invoice | > 99.5% |
| **Data loss rate** | % of ingested events lost in processing | < 0.01% |
| **Pricing catalog freshness** | Max days since last pricing update per provider | < 7 days |
| **Support response time** | Median first response time | < 24 hours (Growth), < 4 hours (Enterprise) |

---

## North Star Metric

**Monthly attributed AI spend tracked on platform** (total dollars flowing through the system).

Why this metric:
- Directly correlates with customer value (more spend tracked = more visibility)
- Correlates with revenue (higher spend customers pay more)
- Correlates with product stickiness (historical data creates switching cost)
- Measurable from day one

| Milestone | Target |
|-----------|--------|
| Month 3 | $1M/month tracked |
| Month 6 | $10M/month tracked |
| Month 12 | $100M/month tracked |
| Month 24 | $1B/month tracked |

---

## Phase-Specific Success Criteria

### Phase 0 (Product Definition) — Current
- [ ] 12 product definition documents completed
- [ ] 5 customer discovery interviews conducted
- [ ] Problem validated by 3+ prospects confirming pain
- [ ] Pricing validated by 3+ prospects confirming willingness to pay

### Phase 0.5 (Technical Architecture)
- [ ] Architecture reviewed by 1+ senior engineer
- [ ] Data model supports all Phase 1 functional requirements
- [ ] Security design passes basic threat model review

### Phase 1 (Ingestion MVP)
- [ ] 3 design partners integrated and viewing dashboards
- [ ] Time to first insight < 24 hours (measured)
- [ ] Cost accuracy > 99% vs. provider invoices (validated)
- [ ] 1 paying customer

### Phase 2 (Budgets + Governance)
- [ ] 10 paying customers
- [ ] Budget alerts deployed by > 50% of customers
- [ ] Waste detection identifying > $10K/month per customer (avg)
- [ ] $25K MRR

### Phase 3 (Optimization + ROI)
- [ ] 30 paying customers
- [ ] Optimization recommendations acted on by > 30% of customers
- [ ] $50K MRR
- [ ] 1 enterprise customer (>$5K/mo)

### Phase 4 (Enterprise)
- [ ] SSO/SCIM deployed
- [ ] SOC 2 Type II in progress
- [ ] 3 enterprise customers
- [ ] $100K MRR

---

## KPI Dashboard (Internal)

The founding team should review weekly:

| Metric | Source |
|--------|--------|
| MRR / ARR | Billing system |
| Active customers | CRM |
| Trial conversions | Product analytics |
| North star (spend tracked) | Platform database |
| Ingest uptime / latency | Monitoring (Datadog/Grafana) |
| Support tickets | Help desk |
| NPS | Quarterly survey |