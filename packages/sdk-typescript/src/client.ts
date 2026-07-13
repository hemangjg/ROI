import type {
  AuthResponse,
  BatchIngestResponse,
  EventStatus,
  IngestEventInput,
  IngestResponse,
  LoginRequest,
  RefreshTokenRequest,
  RegisterRequest,
  UsageEvent,
} from "@ai-finops/api-schemas";
import { HttpClient } from "@ai-finops/sdk-core";

import { normalizeUsageEvent } from "./ingest.js";

export type FinOpsClientConfig = {
  baseUrl: string;
  apiKey?: string;
  accessToken?: string;
  timeoutMs?: number;
  fetch?: typeof fetch;
};

export class FinOpsClient {
  private readonly http: HttpClient;

  constructor(config: FinOpsClientConfig) {
    this.http = new HttpClient({
      baseUrl: config.baseUrl,
      apiKey: config.apiKey,
      accessToken: config.accessToken,
      timeoutMs: config.timeoutMs,
      fetch: config.fetch,
      userAgent: "ai-finops-sdk/0.0.0",
    });
  }

  /** POST /events — ingest a single usage event (API key auth). */
  async ingest(event: IngestEventInput): Promise<IngestResponse> {
    const payload = normalizeUsageEvent(event);
    return this.http.post<IngestResponse>("/events", payload, { auth: "api-key" });
  }

  /** Fire-and-forget alias used by provider instrumentation helpers. */
  track(event: IngestEventInput): void {
    void this.ingest(event).catch(() => {
      // Phase 2 adds a local retry buffer; foundation SDK swallows background errors.
    });
  }

  /** POST /events/batch — ingest up to 500 events (API key auth). */
  async ingestBatch(events: IngestEventInput[]): Promise<BatchIngestResponse> {
    const payload = {
      events: events.map((event) => normalizeUsageEvent(event)),
    };

    return this.http.post<BatchIngestResponse>("/events/batch", payload, { auth: "api-key" });
  }

  /** GET /events/{eventId} — poll processing status (API key auth). */
  async getEventStatus(eventId: string): Promise<EventStatus> {
    return this.http.get<EventStatus>(`/events/${eventId}`, { auth: "api-key" });
  }

  /** POST /auth/login — exchange credentials for JWT pair (no auth header). */
  async login(request: LoginRequest): Promise<AuthResponse> {
    return this.http.post<AuthResponse>("/auth/login", request, { auth: "none" });
  }

  /** POST /auth/register — create org + user (no auth header). */
  async register(request: RegisterRequest): Promise<AuthResponse> {
    return this.http.post<AuthResponse>("/auth/register", request, { auth: "none" });
  }

  /** POST /auth/refresh — rotate JWT pair (no auth header). */
  async refresh(request: RefreshTokenRequest): Promise<AuthResponse> {
    return this.http.post<AuthResponse>("/auth/refresh", request, { auth: "none" });
  }

  /** Convenience helper for strongly typed batch payloads. */
  toUsageEvent(event: IngestEventInput): UsageEvent {
    return normalizeUsageEvent(event);
  }
}
