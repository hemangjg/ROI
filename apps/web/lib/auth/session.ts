export const SESSION_COOKIE = "finops_session";
export const ORG_COOKIE = "finops_org_id";

export type SessionPayload = {
  accessToken: string;
  orgId: string;
};

/** URL-safe base64 without relying on Buffer.base64url (missing in some polyfills). */
function encodeBase64Url(input: string): string {
  let base64: string;
  if (typeof Buffer !== "undefined") {
    base64 = Buffer.from(input, "utf8").toString("base64");
  } else {
    const bytes = new TextEncoder().encode(input);
    let binary = "";
    for (const byte of bytes) {
      binary += String.fromCharCode(byte);
    }
    base64 = btoa(binary);
  }
  return base64.replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/g, "");
}

function decodeBase64Url(input: string): string {
  const base64 = input.replace(/-/g, "+").replace(/_/g, "/");
  const padded = base64 + "=".repeat((4 - (base64.length % 4)) % 4);

  if (typeof Buffer !== "undefined") {
    return Buffer.from(padded, "base64").toString("utf8");
  }

  const binary = atob(padded);
  const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0));
  return new TextDecoder().decode(bytes);
}

export function encodeSession(payload: SessionPayload): string {
  return encodeBase64Url(JSON.stringify(payload));
}

export function decodeSession(value: string): SessionPayload | null {
  try {
    const parsed = JSON.parse(decodeBase64Url(value)) as SessionPayload;
    if (parsed.accessToken && parsed.orgId) return parsed;
    return null;
  } catch {
    return null;
  }
}
