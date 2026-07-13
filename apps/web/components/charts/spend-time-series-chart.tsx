"use client";

import type { SpendTimeSeries } from "@ai-finops/api-schemas";
import ReactECharts from "echarts-for-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

type SpendTimeSeriesChartProps = {
  data?: SpendTimeSeries;
  isLoading: boolean;
};

export function SpendTimeSeriesChart({ data, isLoading }: SpendTimeSeriesChartProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-5 w-40" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-72 w-full" />
        </CardContent>
      </Card>
    );
  }

  const option = {
    grid: { left: 48, right: 24, top: 24, bottom: 40 },
    tooltip: {
      trigger: "axis",
      valueFormatter: (value: number) =>
        new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(value),
    },
    xAxis: {
      type: "category",
      data: data?.data.map((point) => point.date) ?? [],
      axisLine: { lineStyle: { color: "#cbd5e1" } },
    },
    yAxis: {
      type: "value",
      axisLabel: {
        formatter: (value: number) => `$${value}`,
      },
      splitLine: { lineStyle: { color: "#f1f5f9" } },
    },
    series: [
      {
        type: "line",
        smooth: true,
        data: data?.data.map((point) => Number.parseFloat(point.cost_usd)) ?? [],
        lineStyle: { width: 3, color: "#0f172a" },
        areaStyle: { color: "rgba(15, 23, 42, 0.06)" },
        symbol: "circle",
        symbolSize: 6,
      },
    ],
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Daily spend</CardTitle>
      </CardHeader>
      <CardContent>
        <ReactECharts option={option} style={{ height: 320 }} opts={{ renderer: "svg" }} />
      </CardContent>
    </Card>
  );
}
