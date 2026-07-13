"use client";

import type { SpendSummary } from "@ai-finops/api-schemas";
import { ArrowDownRight, ArrowUpRight } from "lucide-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { formatUsd } from "@/lib/utils";

type SpendSummaryCardsProps = {
  data?: SpendSummary;
  isLoading: boolean;
};

export function SpendSummaryCards({ data, isLoading }: SpendSummaryCardsProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <Card key={i}>
            <CardHeader>
              <Skeleton className="h-4 w-24" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-32" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  const change = data?.change_pct ?? 0;
  const ChangeIcon = change >= 0 ? ArrowUpRight : ArrowDownRight;

  return (
    <div className="grid gap-4 md:grid-cols-3">
      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium text-muted-foreground">Today</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-semibold tracking-tight">
            {formatUsd(data?.today_usd ?? "0")}
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium text-muted-foreground">Month to date</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-semibold tracking-tight">{formatUsd(data?.mtd_usd ?? "0")}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium text-muted-foreground">
            vs previous period
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-2">
            <p className="text-3xl font-semibold tracking-tight">{change.toFixed(1)}%</p>
            <ChangeIcon
              className={`h-5 w-5 ${change >= 0 ? "text-emerald-600" : "text-rose-600"}`}
            />
          </div>
          <p className="mt-1 text-sm text-muted-foreground">
            Previous: {formatUsd(data?.previous_period_usd ?? "0")}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
