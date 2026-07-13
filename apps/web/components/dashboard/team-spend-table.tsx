"use client";

import type { SpendByTeam } from "@ai-finops/api-schemas";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { formatUsd } from "@/lib/utils";
import { useDashboardStore } from "@/lib/stores/dashboard-store";

type TeamRow = SpendByTeam["teams"][number];

const columnHelper = createColumnHelper<TeamRow>();

const columns = [
  columnHelper.accessor("name", { header: "Team" }),
  columnHelper.accessor("cost_usd", {
    header: "Spend",
    cell: (info) => formatUsd(info.getValue()),
  }),
  columnHelper.accessor("event_count", {
    header: "Events",
    cell: (info) => info.getValue().toLocaleString(),
  }),
];

type TeamSpendTableProps = {
  data?: SpendByTeam;
  isLoading: boolean;
};

export function TeamSpendTable({ data, isLoading }: TeamSpendTableProps) {
  const setSelectedTeamId = useDashboardStore((s) => s.setSelectedTeamId);

  const table = useReactTable({
    data: data?.teams ?? [],
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-5 w-32" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-48 w-full" />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Team spend</CardTitle>
      </CardHeader>
      <CardContent>
        <table className="w-full text-sm">
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id} className="border-b text-left text-muted-foreground">
                {headerGroup.headers.map((header) => (
                  <th key={header.id} className="pb-3 font-medium">
                    {flexRender(header.column.columnDef.header, header.getContext())}
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {table.getRowModel().rows.map((row) => (
              <tr
                key={row.id}
                className="cursor-pointer border-b transition-colors hover:bg-muted/50"
                onClick={() => setSelectedTeamId(row.original.id)}
              >
                {row.getVisibleCells().map((cell) => (
                  <td key={cell.id} className="py-3">
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </CardContent>
    </Card>
  );
}
