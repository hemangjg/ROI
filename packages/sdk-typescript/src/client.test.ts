import { describe, expect, it, vi } from "vitest";

import { FinOpsClient } from "./client.js";
import { normalizeUsageEvent } from "./ingest.js";

describe("FinOpsClient", () => {
  it("matches the dashboard integration guide example", async () => {
    const fetchMock = vi.fn(
      async () =>
        new Response(
          JSON.stringify({ event_id: "550e8400-e29b-41d4-a716-446655440099", status: "queued" }),
          {
            status: 202,
            headers: { "Content-Type": "application/json" },
          },
        ),
    );

    const client = new FinOpsClient({
      apiKey: "aif_test_key",
      baseUrl: "http://localhost:8888/v1",
      fetch: fetchMock,
    });

    const result = await client.ingest({
      provider: "openai",
      model: "gpt-4o",
      input_tokens: 1200,
      output_tokens: 340,
      occurred_at: new Date().toISOString(),
    });

    expect(result.status).toBe("queued");
    expect(fetchMock).toHaveBeenCalledOnce();

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit];
    const body = JSON.parse(String(init.body)) as Record<string, unknown>;
    expect(body.provider).toBe("openai");
    expect(body.model).toBe("gpt-4o");
    expect(typeof body.idempotency_key).toBe("string");
  });

  it("exposes auth helpers without Authorization header", async () => {
    const fetchMock = vi.fn(
      async () =>
        new Response(
          JSON.stringify({
            access_token: "jwt-access",
            refresh_token: "jwt-refresh",
            org_id: "org-1",
          }),
          {
            status: 200,
            headers: { "Content-Type": "application/json" },
          },
        ),
    );

    const client = new FinOpsClient({
      baseUrl: "http://localhost:8888/v1",
      fetch: fetchMock,
    });

    await client.login({ email: "ops@example.com", password: "secret" });

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit];
    const headers = init.headers as Record<string, string>;
    expect(headers.Authorization).toBeUndefined();
  });
});

describe("normalizeUsageEvent", () => {
  it("fills idempotency_key and cache token defaults", () => {
    const event = normalizeUsageEvent({
      provider: "openai",
      model: "gpt-4o",
      input_tokens: 1,
      output_tokens: 2,
      occurred_at: "2026-06-28T10:00:00Z",
    });

    expect(event.idempotency_key.length).toBeGreaterThan(0);
    expect(event.cache_read_tokens).toBe(0);
    expect(event.cache_write_tokens).toBe(0);
  });
});
