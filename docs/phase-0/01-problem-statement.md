# 01 — Problem Statement

## Executive Summary

Organizations are deploying AI at scale without the financial visibility, governance controls, or optimization tooling that mature cloud infrastructure already has. AI spend is growing faster than any other software line item, yet most companies cannot answer basic questions: How much did we spend? Who spent it? Was it worth it?

This platform exists to become the system of record for enterprise AI spending — the equivalent of what Datadog, CloudZero, or Kubecost did for cloud infrastructure, applied to Large Language Model usage.

---

## The Problem

### 1. AI Spend Is Invisible

Most companies discover their AI costs through monthly provider invoices — after the money is spent. There is no real-time view of consumption by team, project, feature, or customer. Finance receives a single line item ("OpenAI: $47,000") with no attribution.

**Who feels this:** CFO, FinOps Manager, VP Engineering

### 2. AI Spend Is Unallocated

Cloud FinOps solved chargeback and showback a decade ago. AI has no equivalent. Engineering teams experiment freely; product teams embed LLMs in features without cost awareness; individual employees use ChatGPT Enterprise alongside custom API integrations — all without centralized attribution.

**Who feels this:** Finance, Engineering Managers, Product Leaders

### 3. AI Spend Is Ungoverned

Security teams cannot enforce which models are approved, which data can be sent to external providers, or which teams are authorized to use AI at all. Shadow AI — employees using personal accounts or unapproved tools — creates compliance and data leakage risk.

**Who feels this:** CISO, Security Engineering, Legal/Compliance

### 4. AI Spend Is Wasteful

Without visibility, waste compounds silently:

- Duplicate prompts sent to expensive models
- GPT-4 used where GPT-4o-mini would suffice
- Massive context windows loaded unnecessarily
- No caching of repeated queries
- Retry loops burning tokens on failures

Industry estimates suggest 15–30% of LLM spend is recoverable waste. At $500K/year AI spend, that is $75K–$150K left on the table.

**Who feels this:** CTO, Platform Engineering, FinOps

### 5. AI ROI Is Unprovable

Boards and investors now ask: "What is our AI ROI?" Most companies cannot answer. They know what they spent but not what they got — no linkage between AI cost and business outcomes (tickets resolved, features shipped, revenue influenced).

**Who feels this:** CEO, CFO, Board

### 6. Multi-Provider Sprawl

Enterprises rarely use a single LLM provider. A typical mid-market SaaS company runs:

- OpenAI for production features
- Anthropic for internal tooling
- Azure OpenAI for enterprise compliance
- Copilot for developer productivity
- Gemini for specific integrations
- Ollama for on-prem experimentation

Each provider has its own billing console, its own usage format, its own pricing model. No unified view exists.

**Who feels this:** Platform Engineering, FinOps, CTO

---

## Why Now

| Signal | Evidence |
|--------|----------|
| Spend magnitude | Mid-market companies now spend $100K–$2M/year on LLM APIs |
| Growth rate | AI API spend growing 3–10x year-over-year at adopters |
| Board pressure | "AI strategy" and "AI ROI" are standard board agenda items |
| Regulatory pressure | EU AI Act, SOC 2 auditors asking about AI governance |
| Tooling gap | Cloud FinOps is mature; AI FinOps category is nascent |
| Provider fragmentation | 10+ major providers, each with different pricing and APIs |

---

## Job To Be Done

When a company hires this platform, they are hiring it to answer five questions:

1. **How much are we spending on AI?** — Real-time, accurate, cross-provider
2. **Who is spending it?** — Attribution by team, user, project, feature, customer
3. **Are we wasting money?** — Detect and surface optimization opportunities
4. **Are we within budget?** — Enforce limits, alert before overspend
5. **Is it worth it?** — Connect spend to productivity and business outcomes

---

## What This Is Not

| Not this | Why |
|----------|-----|
| A simple cost dashboard | Dashboards without ingestion, attribution, and governance do not retain customers |
| A dev observability tool | LangSmith and Helicone trace requests; they do not do FinOps |
| A single-provider billing viewer | OpenAI's dashboard only shows OpenAI spend |
| A prompt engineering tool | Prompt optimization is a feature, not the product |
| A replacement for LLM providers | We sit alongside providers, not in front of them (initially) |

---

## Problem Validation Criteria

We know this problem is real when:

- A prospect says "I got a $200K OpenAI bill and had no idea until finance flagged it"
- FinOps teams manually export CSVs from 3+ provider consoles and merge in spreadsheets
- Security teams discover employees sending customer PII to unapproved models
- Engineering managers cannot justify AI budget requests to finance
- Board decks include "AI spend" as a line item with a question mark next to ROI

---

## Success Definition for Phase 0

Phase 0 does not build the product. It defines whether the product should be built, for whom, and how it wins. The problem statement above must be validated against target customer interviews before Phase 0.5 (Technical Architecture) begins.