import { describe, expect, it, vi } from "vitest";

import { FinOpsHttpError } from "./errors.js";
import { HttpClient } from "./http-client.js";

describe("HttpClient", () => {
  it("sends bearer auth and JSON body", async () => {
    const fetchMock = vi.fn(async (_input: RequestInfo | URL, init?: RequestInit) => {
      expect(init?.headers).toMatchObject({
        Authorization: "Bearer aif_test_key",
        "Content-Type": "application/json",
      });

      return new Response(JSON.stringify({ event_id: "evt_1", status: "queued" }), {
        status: 202,
        headers: { "Content-Type": "application/json" },
      });
    });

    const client = new HttpClient({
      baseUrl: "http://localhost:8888/v1",
      apiKey: "aif_test_key",
      fetch: fetchMock,
    });

    const result = await client.post("/events", { provider: "openai" });
    expect(result).toEqual({ event_id: "evt_1", status: "queued" });
    expect(fetchMock).toHaveBeenCalledWith(
      "http://localhost:8888/v1/events",
      expect.objectContaining({ method: "POST" }),
    );
  });

  it("throws FinOpsHttpError for non-2xx responses", async () => {
    const fetchMock = vi.fn(
      async () =>
        new Response(
          JSON.stringify({
            title: "Unauthorized",
            status: 401,
            detail: "Invalid API key",
          }),
          {
            status: 401,
            headers: { "Content-Type": "application/problem+json" },
          },
        ),
    );

    const client = new HttpClient({
      baseUrl: "http://localhost:8888/v1",
      apiKey: "aif_bad",
      fetch: fetchMock,
    });

    await expect(client.post("/events", {})).rejects.toMatchObject({
      name: "FinOpsHttpError",
      status: 401,
      message: "Invalid API key",
    } satisfies Partial<FinOpsHttpError>);
  });
});
