import { NextResponse } from "next/server";

import { DEMO_ORG_ID } from "@/lib/api/mock-data";
import { ORG_COOKIE, SESSION_COOKIE, encodeSession } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export async function POST(request: Request) {
  const body = (await request.json()) as { email?: string; password?: string };
  const useMocks = process.env.NEXT_PUBLIC_USE_MOCK_API !== "false";

  let accessToken = "mock-access-token";
  let orgId = DEMO_ORG_ID;

  if (!useMocks) {
    const apiBase = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8888/v1";
    const upstream = await fetch(`${apiBase}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    if (!upstream.ok) {
      return NextResponse.json({ error: "invalid credentials" }, { status: 401 });
    }

    const data = (await upstream.json()) as { access_token: string; org_id: string };
    accessToken = data.access_token;
    orgId = data.org_id;
  }

  const session = encodeSession({ accessToken, orgId });
  const response = NextResponse.json({ ok: true, org_id: orgId });
  response.cookies.set(SESSION_COOKIE, session, { path: "/", sameSite: "lax" });
  response.cookies.set(ORG_COOKIE, orgId, { path: "/", sameSite: "lax" });
  return response;
}
