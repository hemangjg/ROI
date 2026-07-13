"use client";

import { Suspense } from "react";
import { useQuery } from "@tanstack/react-query";

import { ProviderBreakdownBar } from "@/components/charts/provider-breakdown-bar";
import { AppHeader } from "@/components/layout/app-header";
import { useDateRange } from "@/hooks/use-date-range";
import { useOrgId } from "@/hooks/use-org-id";
import { api } from "@/lib/api/client";

function ProvidersContent() {
  const orgId = useOrgId();
  const { range } = useDateRange();
  const providers = useQuery({
    queryKey: ["spend", "by-provider", orgId, range],
    queryFn: () => api.analytics.getSpendByProvider(orgId, range),
    staleTime: 60_000,
  });

  return (
    <main className="px-8 py-6">
      <ProviderBreakdownBar data={providers.data} isLoading={providers.isLoading} />
    </main>
  );
}

export default function ProvidersPage() {
  return (
    <>
      <AppHeader title="Providers" />
      <Suspense fallback={<main className="px-8 py-6">Loading...</main>}>
        <ProvidersContent />
      </Suspense>
    </>
  );
}
