import { ORG_COOKIE, SESSION_COOKIE, decodeSession } from "@/lib/auth/session";

function readCookie(name: string): string | null {
  if (typeof document === "undefined") {
    return null;
  }

  const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : null;
}

export function getClientAccessToken(): string | null {
  const raw = readCookie(SESSION_COOKIE);
  if (!raw) {
    return null;
  }

  const session = decodeSession(raw);
  return session?.accessToken ?? null;
}

export function getClientOrgId(): string | null {
  return readCookie(ORG_COOKIE);
}
