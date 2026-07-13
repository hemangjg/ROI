"use client";

import { useQuery } from "@tanstack/react-query";

import type { DateRange } from "@ai-finops/api-schemas";
import { api } from "@/lib/api/client";

export function useSpendSummary(orgId: string, range: DateRange) {
  return useQuery({
    queryKey: ["spend", "summary", orgId, range],
    queryFn: () => api.analytics.getSpendSummary(orgId, range),
    staleTime: 60_000,
  });
}
