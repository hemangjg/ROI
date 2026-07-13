import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AppHeader } from "@/components/layout/app-header";

export default function SettingsTeamsPage() {
  return (
    <>
      <AppHeader title="Team management" showDateRange={false} />
      <main className="px-8 py-6">
        <Card>
          <CardHeader>
            <CardTitle>Teams</CardTitle>
            <CardDescription>Team CRUD via Management API lands in Phase 2.</CardDescription>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground">
            Use the dashboard Teams view for spend breakdown while admin mutations are being built.
          </CardContent>
        </Card>
      </main>
    </>
  );
}
