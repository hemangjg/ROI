"use client";

import { Suspense } from "react";
import { useQuery } from "@tanstack/react-query";

import { TeamSpendTable } from "@/components/dashboard/team-spend-table";
import { AppHeader } from "@/components/layout/app-header";
import { useDateRange } from "@/hooks/use-date-range";
import { useOrgId } from "@/hooks/use-org-id";
import { api } from "@/lib/api/client";

function TeamsContent() {
  const orgId = useOrgId();
  const { range } = useDateRange();
  const teams = useQuery({
    queryKey: ["spend", "by-team", orgId, range],
    queryFn: () => api.analytics.getSpendByTeam(orgId, range),
    staleTime: 60_000,
  });

  return (
    <main className="px-8 py-6">
      <TeamSpendTable data={teams.data} isLoading={teams.isLoading} />
    </main>
  );
}

export default function TeamsPage() {
  return (
    <>
      <AppHeader title="Teams" />
      <Suspense fallback={<main className="px-8 py-6">Loading...</main>}>
        <TeamsContent />
      </Suspense>
    </>
  );
}
