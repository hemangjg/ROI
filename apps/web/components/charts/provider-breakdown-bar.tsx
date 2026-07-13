"use client";

import type { SpendByProvider } from "@ai-finops/api-schemas";
import ReactECharts from "echarts-for-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

const PROVIDER_COLORS: Record<string, string> = {
  openai: "#10a37f",
  anthropic: "#d97757",
  google: "#4285f4",
};

type ProviderBreakdownBarProps = {
  data?: SpendByProvider;
  isLoading: boolean;
};

export function ProviderBreakdownBar({ data, isLoading }: ProviderBreakdownBarProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-5 w-48" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-72 w-full" />
        </CardContent>
      </Card>
    );
  }

  const providers = data?.providers ?? [];
  const option = {
    grid: { left: 120, right: 24, top: 24, bottom: 24 },
    tooltip: {
      trigger: "axis",
      axisPointer: { type: "shadow" },
      valueFormatter: (value: number) =>
        new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(value),
    },
    xAxis: { type: "value", splitLine: { lineStyle: { color: "#f1f5f9" } } },
    yAxis: {
      type: "category",
      data: providers.map((p) => p.name),
      axisLine: { show: false },
      axisTick: { show: false },
    },
    series: [
      {
        type: "bar",
        data: providers.map((p) => ({
          value: Number.parseFloat(p.cost_usd),
          itemStyle: { color: PROVIDER_COLORS[p.name] ?? "#64748b", borderRadius: [0, 4, 4, 0] },
        })),
        barWidth: 24,
      },
    ],
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Spend by provider</CardTitle>
      </CardHeader>
      <CardContent>
        <ReactECharts option={option} style={{ height: 280 }} opts={{ renderer: "svg" }} />
      </CardContent>
    </Card>
  );
}
