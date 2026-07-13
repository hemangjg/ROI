"use client";

import { Suspense } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import { DateRangePicker } from "@/components/dashboard/date-range-picker";

export function AppHeader({
  title,
  showDateRange = true,
}: {
  title: string;
  showDateRange?: boolean;
}) {
  const router = useRouter();

  return (
    <header className="flex items-center justify-between border-b bg-background px-8 py-5">
      <div>
        <p className="text-sm text-muted-foreground">Organization spend</p>
        <h2 className="text-2xl font-semibold tracking-tight">{title}</h2>
      </div>
      <div className="flex items-center gap-4">
        {showDateRange ? (
          <Suspense fallback={<div className="h-9 w-40 animate-pulse rounded-md bg-muted" />}>
            <DateRangePicker />
          </Suspense>
        ) : null}
        <Button
          variant="outline"
          size="sm"
          onClick={() => {
            document.cookie = "finops_session=; Max-Age=0; path=/";
            document.cookie = "finops_org_id=; Max-Age=0; path=/";
            router.push("/login");
          }}
        >
          Sign out
        </Button>
      </div>
    </header>
  );
}
