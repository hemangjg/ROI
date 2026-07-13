"use client";

import { Suspense } from "react";
import { useQuery } from "@tanstack/react-query";

import { SpendTimeSeriesChart } from "@/components/charts/spend-time-series-chart";
import { EmptyStateIntegrationGuide } from "@/components/dashboard/empty-state-integration-guide";
import { SpendSummaryCards } from "@/components/dashboard/spend-summary-cards";
import { AppHeader } from "@/components/layout/app-header";
import { useDateRange } from "@/hooks/use-date-range";
import { useOrgId } from "@/hooks/use-org-id";
import { useSpendSummary } from "@/hooks/use-spend-summary";
import { api } from "@/lib/api/client";

function DashboardContent() {
  const orgId = useOrgId();
  const { range } = useDateRange();
  const summary = useSpendSummary(orgId, range);
  const timeseries = useQuery({
    queryKey: ["spend", "timeseries", orgId, range],
    queryFn: () => api.analytics.getSpendTimeSeries(orgId, range),
    staleTime: 60_000,
  });

  return (
    <main className="space-y-6 px-8 py-6">
      <SpendSummaryCards data={summary.data} isLoading={summary.isLoading} />
      <SpendTimeSeriesChart data={timeseries.data} isLoading={timeseries.isLoading} />
      <EmptyStateIntegrationGuide />
    </main>
  );
}

export default function DashboardPage() {
  return (
    <>
      <AppHeader title="Overview" />
      <Suspense fallback={<main className="px-8 py-6">Loading dashboard...</main>}>
        <DashboardContent />
      </Suspense>
    </>
  );
}
