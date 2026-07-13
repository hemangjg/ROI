"use client";

import { useEffect, useState } from "react";

import { DEMO_ORG_ID } from "@/lib/api/mock-data";
import { getClientOrgId } from "@/lib/auth/client-session";

export function useOrgId(): string {
  const [orgId, setOrgId] = useState(DEMO_ORG_ID);

  useEffect(() => {
    const fromCookie = getClientOrgId();
    if (fromCookie) {
      setOrgId(fromCookie);
    }
  }, []);

  return orgId;
}
