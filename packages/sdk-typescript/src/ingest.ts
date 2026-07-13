import type { IngestEventInput, UsageEvent } from "@ai-finops/api-schemas";

export function createIdempotencyKey(): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID();
  }

  return `evt_${Date.now()}_${Math.random().toString(36).slice(2, 10)}`;
}

export function normalizeUsageEvent(event: IngestEventInput): UsageEvent {
  return {
    ...event,
    idempotency_key: event.idempotency_key ?? createIdempotencyKey(),
    cache_read_tokens: event.cache_read_tokens ?? 0,
    cache_write_tokens: event.cache_write_tokens ?? 0,
  };
}
