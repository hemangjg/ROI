"use client";

import { useQuery } from "@tanstack/react-query";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useOrgId } from "@/hooks/use-org-id";
import { api } from "@/lib/api/client";

export function ApiKeyManager() {
  const orgId = useOrgId();
  const keys = useQuery({
    queryKey: ["management", "api-keys", orgId],
    queryFn: () => api.management.listApiKeys(orgId),
    staleTime: 60_000,
  });

  if (keys.isLoading) {
    return <Skeleton className="h-48 w-full" />;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>API keys</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {(keys.data ?? []).map((key) => (
          <div key={key.id} className="flex items-center justify-between rounded-lg border p-4">
            <div>
              <p className="font-medium">{key.name}</p>
              <p className="text-sm text-muted-foreground">{key.key_prefix}••••••••</p>
            </div>
            <Badge variant="secondary">Active</Badge>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
