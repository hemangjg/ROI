"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import {
  Building2,
  KeyRound,
  LayoutDashboard,
  LineChart,
  ScrollText,
  Settings,
  Users,
  Wallet,
} from "lucide-react";

import { useOrgId } from "@/hooks/use-org-id";
import { api } from "@/lib/api/client";
import { cn } from "@/lib/utils";

const navItems = [
  { href: "/dashboard", label: "Overview", icon: LayoutDashboard },
  { href: "/dashboard/providers", label: "Providers", icon: LineChart },
  { href: "/dashboard/teams", label: "Teams", icon: Users },
  { href: "/dashboard/budgets", label: "Budgets", icon: Wallet, badgeKey: "alerts" as const },
];

const settingsItems = [
  { href: "/settings/api-keys", label: "API keys", icon: KeyRound },
  { href: "/settings/organization", label: "Organization", icon: Building2 },
  { href: "/settings/teams", label: "Teams", icon: Users },
  { href: "/settings/audit", label: "Audit log", icon: ScrollText },
];

export function AppSidebar() {
  const pathname = usePathname();
  const orgId = useOrgId();
  const alerts = useQuery({
    queryKey: ["budget-alerts", "open", orgId],
    queryFn: () => api.management.listBudgetAlerts(orgId, "open"),
    staleTime: 30_000,
  });
  const openAlerts = alerts.data?.length ?? 0;

  return (
    <aside className="flex h-full w-64 flex-col border-r bg-card">
      <div className="border-b px-6 py-5">
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          AI FinOps
        </p>
        <h1 className="mt-1 text-lg font-semibold tracking-tight">Spend Control</h1>
      </div>
      <nav className="flex-1 space-y-6 px-4 py-6">
        <div>
          <p className="mb-2 px-2 text-xs font-medium uppercase tracking-wide text-muted-foreground">
            Dashboard
          </p>
          <div className="space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const active = pathname === item.href;
              const showBadge = item.badgeKey === "alerts" && openAlerts > 0;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors",
                    active
                      ? "bg-primary text-primary-foreground"
                      : "text-foreground hover:bg-muted",
                  )}
                >
                  <Icon className="h-4 w-4" />
                  <span className="flex-1">{item.label}</span>
                  {showBadge ? (
                    <span
                      className={cn(
                        "rounded-full px-1.5 py-0.5 text-[10px] font-semibold",
                        active
                          ? "bg-primary-foreground/20 text-primary-foreground"
                          : "bg-destructive text-destructive-foreground",
                      )}
                    >
                      {openAlerts}
                    </span>
                  ) : null}
                </Link>
              );
            })}
          </div>
        </div>
        <div>
          <p className="mb-2 flex items-center gap-2 px-2 text-xs font-medium uppercase tracking-wide text-muted-foreground">
            <Settings className="h-3 w-3" />
            Settings
          </p>
          <div className="space-y-1">
            {settingsItems.map((item) => {
              const Icon = item.icon;
              const active = pathname === item.href;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors",
                    active
                      ? "bg-primary text-primary-foreground"
                      : "text-foreground hover:bg-muted",
                  )}
                >
                  <Icon className="h-4 w-4" />
                  {item.label}
                </Link>
              );
            })}
          </div>
        </div>
      </nav>
    </aside>
  );
}
