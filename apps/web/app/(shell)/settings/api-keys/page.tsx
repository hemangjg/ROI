import { ApiKeyManager } from "@/components/settings/api-key-manager";
import { AppHeader } from "@/components/layout/app-header";

export default function ApiKeysPage() {
  return (
    <>
      <AppHeader title="API keys" showDateRange={false} />
      <main className="px-8 py-6">
        <ApiKeyManager />
      </main>
    </>
  );
}
