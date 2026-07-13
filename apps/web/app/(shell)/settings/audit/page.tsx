import { AuditLogViewer } from "@/components/settings/audit-log-viewer";
import { AppHeader } from "@/components/layout/app-header";

export default function AuditPage() {
  return (
    <>
      <AppHeader title="Audit log" showDateRange={false} />
      <main className="px-8 py-6">
        <AuditLogViewer />
      </main>
    </>
  );
}
