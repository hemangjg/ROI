# 06 — Frontend Architecture

## Stack (Locked)

| Concern | Technology |
|---------|-----------|
| Framework | Next.js 15 (App Router) |
| UI | React 19 + shadcn/ui + Tailwind CSS |
| Tables | TanStack Table |
| Charts | Apache ECharts |
| Forms | React Hook Form + Zod |
| State | Zustand |
| Data fetching | TanStack Query |
| Icons | Lucide |

---

## Application Structure

```
apps/web/
├── app/
│   ├── (auth)/
│   │   ├── login/page.tsx
│   │   └── register/page.tsx
│   ├── (dashboard)/
│   │   ├── layout.tsx              # Sidebar + header
│   │   ├── page.tsx                # Spend overview
│   │   ├── providers/page.tsx
│   │   └── teams/page.tsx
│   ├── (settings)/
│   │   ├── api-keys/page.tsx
│   │   ├── organization/page.tsx
│   │   ├── teams/page.tsx
│   │   └── audit/page.tsx
│   └── layout.tsx
├── components/
│   ├── charts/
│   │   ├── SpendTimeSeriesChart.tsx   # ECharts
│   │   └── ProviderBreakdownBar.tsx
│   ├── dashboard/
│   │   ├── SpendSummaryCards.tsx
│   │   └── TeamSpendTable.tsx         # TanStack Table
│   ├── settings/
│   │   ├── ApiKeyManager.tsx
│   │   └── AuditLogViewer.tsx
│   └── ui/                            # shadcn components
├── lib/
│   ├── api/                           # Typed API client
│   ├── auth/                          # Token management
│   └── stores/                        # Zustand stores
└── hooks/
    ├── useSpendSummary.ts
    └── useDateRange.ts
```

---

## Pages (Phase 1)

| Route | Purpose | API calls |
|-------|---------|-----------|
| `/login` | Email/password login | POST `/v1/auth/login` |
| `/register` | Create account + org | POST `/v1/auth/register` |
| `/dashboard` | Spend overview | GET summary + timeseries |
| `/dashboard/providers` | Provider breakdown | GET by-provider |
| `/dashboard/teams` | Team spend table | GET by-team |
| `/settings/api-keys` | Key management | Management API |
| `/settings/organization` | Org profile | Management API |
| `/settings/teams` | Team CRUD | Management API |
| `/settings/audit` | Audit log viewer | GET audit-logs |

---

## Key Components

### SpendSummaryCards

Displays: Today spend, MTD spend, vs previous period (% change).

Uses TanStack Query with 60s stale time. Skeleton loading via shadcn.

### SpendTimeSeriesChart

Apache ECharts line chart. Daily granularity. Responsive. Tooltip shows formatted USD.

### TeamSpendTable

TanStack Table with sorting, pagination. Click row → filter dashboard to team (Zustand `selectedTeamId`).

### DateRangePicker

Presets: 7d, 30d, 90d, custom. Stored in URL search params for shareable links.

### EmptyStateIntegrationGuide

Shown when no events exist. Copy-paste SDK snippet with org API key placeholder.

---

## Data Fetching Pattern

```typescript
// hooks/useSpendSummary.ts
export function useSpendSummary(orgId: string, range: DateRange) {
  return useQuery({
    queryKey: ['spend', 'summary', orgId, range],
    queryFn: () => api.analytics.getSpendSummary(orgId, range),
    staleTime: 60_000,
  });
}
```

- Server Components for layout shell and auth guard
- Client Components for interactive charts and filters
- TanStack Query for all API data with optimistic updates on admin mutations

---

## Auth Flow

1. Login → store JWT in httpOnly cookie (set by API or Next.js middleware)
2. Middleware checks cookie on `(dashboard)` and `(settings)` route groups
3. Refresh token rotation via `/v1/auth/refresh` before access token expiry
4. Logout clears cookies and revokes refresh token

---

## API Client

Typed client generated from `api/openapi/v1.yaml` (orval or openapi-typescript).

Base URL from `NEXT_PUBLIC_API_URL` (Envoy Gateway: `api.ai-finops.local`).

---

## Design Tokens

shadcn/ui default theme with FinOps-specific adjustments:
- Primary: neutral/slate (enterprise, not playful)
- Charts: distinct palette per provider (OpenAI green, Anthropic orange, etc.)
- Typography: Inter (shadcn default)

See design-taste skill principles in Phase 2 UI polish pass.