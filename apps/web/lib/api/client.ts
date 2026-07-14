import type {
  AuditLogPage,
  ApiKey,
  AuthResponse,
  DateRange,
  LoginRequest,
  RegisterRequest,
  SpendByProvider,
  SpendByTeam,
  SpendSummary,
  SpendTimeSeries,
} from "@ai-finops/api-schemas";

import { getClientAccessToken } from "@/lib/auth/client-session";
import {
  DEMO_ORG_ID,
  mockApiKeys,
  mockAuditLog,
  mockSpendByProvider,
  mockSpendByTeam,
  mockSpendSummary,
  mockSpendTimeSeries,
} from "./mock-data";

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8888/v1";

function withDateRange(path: string, range: DateRange): string {
  const params = new URLSearchParams({ from: range.from, to: range.to });
  return `${path}?${params.toString()}`;
}

/** Returns mock data until NEXT_PUBLIC_USE_MOCK_API=false. */
function shouldUseMocks() {
  return process.env.NEXT_PUBLIC_USE_MOCK_API !== "false";
}

function authHeaders(): HeadersInit {
  const token = getClientAccessToken();
  if (!token) {
    return {};
  }
  return { Authorization: `Bearer ${token}` };
}

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers);
  headers.set("Content-Type", "application/json");

  for (const [key, value] of Object.entries(authHeaders())) {
    headers.set(key, value);
  }

  const response = await fetch(`${API_BASE}${path}`, { ...init, headers });
  if (!response.ok) {
    throw new Error(`API ${response.status}: ${path}`);
  }
  return response.json() as Promise<T>;
}

export const api = {
  auth: {
    async login(body: LoginRequest): Promise<AuthResponse> {
      if (shouldUseMocks()) {
        return {
          access_token: "mock-access-token",
          refresh_token: "mock-refresh-token",
          org_id: DEMO_ORG_ID,
        };
      }
      return apiFetch<AuthResponse>("/auth/login", {
        method: "POST",
        body: JSON.stringify(body),
      });
    },
    async register(body: RegisterRequest): Promise<AuthResponse> {
      if (shouldUseMocks()) {
        return {
          access_token: "mock-access-token",
          refresh_token: "mock-refresh-token",
          org_id: DEMO_ORG_ID,
        };
      }
      return apiFetch<AuthResponse>("/auth/register", {
        method: "POST",
        body: JSON.stringify(body),
      });
    },
  },
  analytics: {
    async getSpendSummary(orgId: string, range: DateRange): Promise<SpendSummary> {
      if (shouldUseMocks()) return mockSpendSummary;
      return apiFetch<SpendSummary>(withDateRange(`/orgs/${orgId}/spend/summary`, range));
    },
    async getSpendTimeSeries(orgId: string, range: DateRange): Promise<SpendTimeSeries> {
      if (shouldUseMocks()) return mockSpendTimeSeries;
      return apiFetch<SpendTimeSeries>(withDateRange(`/orgs/${orgId}/spend/timeseries`, range));
    },
    async getSpendByProvider(orgId: string, range: DateRange): Promise<SpendByProvider> {
      if (shouldUseMocks()) return mockSpendByProvider;
      return apiFetch<SpendByProvider>(withDateRange(`/orgs/${orgId}/spend/by-provider`, range));
    },
    async getSpendByTeam(orgId: string, range: DateRange): Promise<SpendByTeam> {
      if (shouldUseMocks()) return mockSpendByTeam;
      return apiFetch<SpendByTeam>(withDateRange(`/orgs/${orgId}/spend/by-team`, range));
    },
  },
  management: {
    async listApiKeys(orgId: string): Promise<ApiKey[]> {
      if (shouldUseMocks()) return mockApiKeys;
      return apiFetch<ApiKey[]>(`/orgs/${orgId}/api-keys`);
    },
    async listAuditLogs(orgId: string): Promise<AuditLogPage> {
      if (shouldUseMocks()) return mockAuditLog;
      return apiFetch<AuditLogPage>(`/orgs/${orgId}/audit-logs`);
    },
    async listTeams(orgId: string): Promise<Team[]> {
      if (shouldUseMocks())
        return [{ id: "team-demo", name: "Platform", created_at: new Date().toISOString() }];
      return apiFetch<Team[]>(`/orgs/${orgId}/teams`);
    },
    async listBudgets(orgId: string, period?: string): Promise<Budget[]> {
      if (shouldUseMocks()) return [];
      const q = period ? `?period=${encodeURIComponent(period)}` : "";
      return apiFetch<Budget[]>(`/orgs/${orgId}/budgets${q}`);
    },
    async upsertBudget(
      orgId: string,
      body: {
        team_id: string;
        amount_usd: string;
        soft_threshold_pct?: number;
        period_start?: string;
      },
    ): Promise<Budget> {
      if (shouldUseMocks()) {
        return {
          id: "budget-demo",
          org_id: orgId,
          team_id: body.team_id,
          amount_usd: body.amount_usd,
          soft_threshold_pct: body.soft_threshold_pct ?? 80,
          period_start: body.period_start ?? "2026-07-01",
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        };
      }
      return apiFetch<Budget>(`/orgs/${orgId}/budgets`, {
        method: "PUT",
        body: JSON.stringify(body),
      });
    },
    async evaluateBudget(
      orgId: string,
      body: { team_id: string; spend_usd: string; period_start?: string },
    ): Promise<BudgetStatus> {
      if (shouldUseMocks()) {
        return {
          id: "budget-demo",
          org_id: orgId,
          team_id: body.team_id,
          amount_usd: "100.00",
          soft_threshold_pct: 80,
          period_start: "2026-07-01",
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          spend_usd: body.spend_usd,
          usage_pct: 85,
          over_soft_threshold: true,
          alert_fired: true,
        };
      }
      return apiFetch<BudgetStatus>(`/orgs/${orgId}/budgets/evaluate`, {
        method: "POST",
        body: JSON.stringify(body),
      });
    },
    async listBudgetAlerts(orgId: string, status?: string): Promise<BudgetAlert[]> {
      if (shouldUseMocks()) return [];
      const q = status ? `?status=${encodeURIComponent(status)}` : "";
      return apiFetch<BudgetAlert[]>(`/orgs/${orgId}/budget-alerts${q}`);
    },
    async acknowledgeBudgetAlert(orgId: string, alertId: string): Promise<BudgetAlert> {
      if (shouldUseMocks()) {
        return {
          id: alertId,
          budget_id: "budget-demo",
          org_id: orgId,
          team_id: "team-demo",
          threshold_pct: 80,
          spend_usd: "85.00",
          budget_usd: "100.00",
          status: "acknowledged",
          message: "mock",
          created_at: new Date().toISOString(),
        };
      }
      return apiFetch<BudgetAlert>(`/orgs/${orgId}/budget-alerts/${alertId}/ack`, {
        method: "POST",
      });
    },
  },
};

export type Team = { id: string; name: string; created_at?: string; org_id?: string };
export type Budget = {
  id: string;
  org_id: string;
  team_id: string;
  team_name?: string;
  amount_usd: string;
  soft_threshold_pct: number;
  period_start: string;
  created_at: string;
  updated_at: string;
};
export type BudgetAlert = {
  id: string;
  budget_id: string;
  org_id: string;
  team_id: string;
  threshold_pct: number;
  spend_usd: string;
  budget_usd: string;
  status: string;
  message: string;
  created_at: string;
};
export type BudgetStatus = Budget & {
  spend_usd: string;
  usage_pct: number;
  over_soft_threshold: boolean;
  alert_fired: boolean;
  alert?: BudgetAlert;
};
