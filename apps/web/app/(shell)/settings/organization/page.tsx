import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AppHeader } from "@/components/layout/app-header";

export default function OrganizationPage() {
  return (
    <>
      <AppHeader title="Organization" showDateRange={false} />
      <main className="px-8 py-6">
        <Card>
          <CardHeader>
            <CardTitle>Demo organization</CardTitle>
            <CardDescription>Organization profile editing ships in Phase 2.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2 text-sm text-muted-foreground">
            <p>Slug: demo</p>
            <p>Plan: Foundation (Phase 1)</p>
          </CardContent>
        </Card>
      </main>
    </>
  );
}
