"use client";

import Link from "next/link";
import { Suspense } from "react";
import { useQuery } from "@tanstack/react-query";
import { AlertTriangle } from "lucide-react";

import { SpendTimeSeriesChart } from "@/components/charts/spend-time-series-chart";
import { EmptyStateIntegrationGuide } from "@/components/dashboard/empty-state-integration-guide";
import { SpendSummaryCards } from "@/components/dashboard/spend-summary-cards";
import { AppHeader } from "@/components/layout/app-header";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
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
  const alerts = useQuery({
    queryKey: ["budget-alerts", "open", orgId],
    queryFn: () => api.management.listBudgetAlerts(orgId, "open"),
    staleTime: 30_000,
  });
  const openCount = alerts.data?.length ?? 0;

  return (
    <main className="space-y-6 px-8 py-6">
      {openCount > 0 ? (
        <Card className="border-amber-500/40 bg-amber-500/5">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <div>
              <CardTitle className="flex items-center gap-2 text-base">
                <AlertTriangle className="h-4 w-4 text-amber-600" />
                {openCount} open budget alert{openCount === 1 ? "" : "s"}
              </CardTitle>
              <CardDescription>Soft limits reached — review spend vs budget.</CardDescription>
            </div>
            <Button asChild size="sm" variant="outline">
              <Link href="/dashboard/budgets">Review budgets</Link>
            </Button>
          </CardHeader>
          <CardContent className="space-y-1 text-sm text-muted-foreground">
            {(alerts.data ?? []).slice(0, 3).map((a) => (
              <p key={a.id} className="truncate">
                {a.message}
              </p>
            ))}
          </CardContent>
        </Card>
      ) : null}
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
