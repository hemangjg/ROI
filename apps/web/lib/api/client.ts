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
  },
};
