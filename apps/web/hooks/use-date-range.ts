"use client";

import { useMemo } from "react";
import { useSearchParams } from "next/navigation";

import type { DateRange } from "@ai-finops/api-schemas";

const PRESETS = {
  "7d": 7,
  "30d": 30,
  "90d": 90,
} as const;

export type DateRangePreset = keyof typeof PRESETS;

function toIsoDate(date: Date) {
  return date.toISOString();
}

export function useDateRange() {
  const searchParams = useSearchParams();
  const preset = (searchParams.get("range") as DateRangePreset | null) ?? "30d";

  const range = useMemo<DateRange>(() => {
    const days = PRESETS[preset] ?? PRESETS["30d"];
    const to = new Date();
    const from = new Date();
    from.setDate(to.getDate() - days);
    return { from: toIsoDate(from), to: toIsoDate(to) };
  }, [preset]);

  return { preset, range };
}
