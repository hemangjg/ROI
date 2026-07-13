/** Hand-written types mirroring api/openapi/v1.yaml. */

export type LoginRequest = {
  email: string;
  password: string;
};

export type RegisterRequest = {
  email: string;
  password: string;
  name: string;
  org_name: string;
};

export type AuthResponse = {
  access_token: string;
  refresh_token: string;
  org_id: string;
};

export type SpendSummary = {
  today_usd: string;
  mtd_usd: string;
  previous_period_usd: string;
  change_pct: number;
};

export type SpendTimeSeriesPoint = {
  date: string;
  cost_usd: string;
};

export type SpendTimeSeries = {
  data: SpendTimeSeriesPoint[];
};

export type SpendByProviderItem = {
  name: string;
  cost_usd: string;
  pct: number;
};

export type SpendByProvider = {
  providers: SpendByProviderItem[];
};

export type SpendByTeamItem = {
  id: string;
  name: string;
  cost_usd: string;
  event_count: number;
  pct?: number;
};

export type SpendByTeam = {
  teams: SpendByTeamItem[];
};

export type ApiKey = {
  id: string;
  name: string;
  key_prefix: string;
  created_at: string;
  last_used_at?: string;
};

export type AuditLogEntry = {
  id: string;
  actor_id?: string;
  actor_type: string;
  action: string;
  resource_type: string;
  resource_id?: string;
  created_at: string;
};

export type AuditLogPage = {
  data: AuditLogEntry[];
  next_cursor?: string;
};

export type DateRange = {
  from: string;
  to: string;
};

export type EventMetadata = {
  team_id?: string;
  user_id?: string;
  project?: string;
  feature?: string;
  customer?: string;
};

export type UsageEvent = {
  idempotency_key: string;
  provider: string;
  model: string;
  input_tokens: number;
  output_tokens: number;
  cache_read_tokens?: number;
  cache_write_tokens?: number;
  occurred_at: string;
  metadata?: EventMetadata;
};

/** Caller may omit idempotency_key; the SDK generates one when missing. */
export type IngestEventInput = Omit<UsageEvent, "idempotency_key"> & {
  idempotency_key?: string;
};

export type IngestResponse = {
  event_id: string;
  status: "queued";
};

export type BatchIngestRequest = {
  events: UsageEvent[];
};

export type BatchIngestError = {
  index: number;
  code: string;
  message: string;
};

export type BatchIngestResponse = {
  accepted: number;
  rejected: number;
  errors?: BatchIngestError[];
};

export type EventStatus = {
  event_id: string;
  status: "queued" | "processed" | "failed";
  cost_usd?: string;
};

export type RefreshTokenRequest = {
  refresh_token: string;
};

export type ApiProblem = {
  type?: string;
  title?: string;
  status?: number;
  detail?: string;
  instance?: string;
};
