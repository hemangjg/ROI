"use client";

import Link from "next/link";
import { usePathname, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";

const PRESETS = [
  { label: "7d", value: "7d" },
  { label: "30d", value: "30d" },
  { label: "90d", value: "90d" },
] as const;

export function DateRangePicker() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const active = searchParams.get("range") ?? "30d";

  return (
    <div className="flex items-center gap-2">
      {PRESETS.map((preset) => {
        const params = new URLSearchParams(searchParams.toString());
        params.set("range", preset.value);
        const href = `${pathname}?${params.toString()}`;
        return (
          <Button
            key={preset.value}
            variant={active === preset.value ? "default" : "outline"}
            size="sm"
            asChild
          >
            <Link href={href}>{preset.label}</Link>
          </Button>
        );
      })}
    </div>
  );
}
