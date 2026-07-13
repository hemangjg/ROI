import type {
  AuditLogPage,
  ApiKey,
  SpendByProvider,
  SpendByTeam,
  SpendSummary,
  SpendTimeSeries,
} from "@ai-finops/api-schemas";

export const DEMO_ORG_ID = "00000000-0000-4000-8000-000000000001";

export const mockSpendSummary: SpendSummary = {
  today_usd: "142.38",
  mtd_usd: "3847.92",
  previous_period_usd: "3291.44",
  change_pct: 16.9,
};

export const mockSpendTimeSeries: SpendTimeSeries = {
  data: [
    { date: "2026-06-24", cost_usd: "412.10" },
    { date: "2026-06-25", cost_usd: "498.22" },
    { date: "2026-06-26", cost_usd: "521.04" },
    { date: "2026-06-27", cost_usd: "467.88" },
    { date: "2026-06-28", cost_usd: "589.31" },
    { date: "2026-06-29", cost_usd: "612.45" },
    { date: "2026-06-30", cost_usd: "746.92" },
  ],
};

export const mockSpendByProvider: SpendByProvider = {
  providers: [
    { name: "openai", cost_usd: "2148.40", pct: 55.8 },
    { name: "anthropic", cost_usd: "1289.12", pct: 33.5 },
    { name: "google", cost_usd: "410.40", pct: 10.7 },
  ],
};

export const mockSpendByTeam: SpendByTeam = {
  teams: [
    { id: "1", name: "Platform", cost_usd: "1420.00", event_count: 420, pct: 36.9 },
    { id: "2", name: "Product", cost_usd: "1180.50", event_count: 318, pct: 30.7 },
    { id: "3", name: "Research", cost_usd: "847.42", event_count: 201, pct: 22.0 },
    { id: "4", name: "Support", cost_usd: "400.00", event_count: 96, pct: 10.4 },
  ],
};

export const mockApiKeys: ApiKey[] = [
  {
    id: "key-1",
    name: "Production ingest",
    key_prefix: "aif_live_",
    created_at: "2026-06-01T10:00:00Z",
    last_used_at: "2026-06-30T18:42:00Z",
  },
  {
    id: "key-2",
    name: "Staging",
    key_prefix: "aif_test_",
    created_at: "2026-06-15T14:30:00Z",
  },
];

export const mockAuditLog: AuditLogPage = {
  data: [
    {
      id: "audit-1",
      actor_type: "user",
      action: "api_key.created",
      resource_type: "api_key",
      resource_id: "key-2",
      created_at: "2026-06-15T14:30:00Z",
    },
    {
      id: "audit-2",
      actor_type: "user",
      action: "team.updated",
      resource_type: "team",
      resource_id: "2",
      created_at: "2026-06-20T09:15:00Z",
    },
  ],
};
