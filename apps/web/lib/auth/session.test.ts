import { describe, expect, it } from "vitest";

import { decodeSession, encodeSession } from "./session";

describe("session encoding", () => {
  const payload = {
    accessToken: "mock-access-token",
    orgId: "00000000-0000-4000-8000-000000000001",
  };

  it("round-trips with Buffer", () => {
    const encoded = encodeSession(payload);
    expect(decodeSession(encoded)).toEqual(payload);
  });

  it("round-trips without Buffer (browser)", () => {
    const originalBuffer = globalThis.Buffer;
    // @ts-expect-error test-only removal
    delete globalThis.Buffer;

    const encoded = encodeSession(payload);
    expect(encoded.length).toBeGreaterThan(0);
    expect(decodeSession(encoded)).toEqual(payload);

    globalThis.Buffer = originalBuffer;
  });

  it("uses standard base64 encoding (not base64url) for compatibility", () => {
    const encoded = encodeSession(payload);
    expect(encoded).not.toContain("+");
    expect(encoded).not.toContain("/");
    expect(encoded).not.toContain("=");
  });
});
