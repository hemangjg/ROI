"use client";

import { useQuery } from "@tanstack/react-query";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useOrgId } from "@/hooks/use-org-id";
import { api } from "@/lib/api/client";

export function AuditLogViewer() {
  const orgId = useOrgId();
  const audit = useQuery({
    queryKey: ["management", "audit-logs", orgId],
    queryFn: () => api.management.listAuditLogs(orgId),
    staleTime: 60_000,
  });

  if (audit.isLoading) {
    return <Skeleton className="h-48 w-full" />;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Audit log</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {(audit.data?.data ?? []).map((entry) => (
            <div key={entry.id} className="rounded-lg border p-4 text-sm">
              <div className="flex items-center justify-between gap-4">
                <p className="font-medium">{entry.action}</p>
                <p className="text-muted-foreground">
                  {new Date(entry.created_at).toLocaleString()}
                </p>
              </div>
              <p className="mt-1 text-muted-foreground">
                {entry.actor_type}
                {entry.actor_id ? `/${entry.actor_id}` : ""} · {entry.resource_type}
                {entry.resource_id ? `/${entry.resource_id}` : ""}
              </p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
