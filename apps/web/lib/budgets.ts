/** MTD date range in UTC (ISO date strings for analytics query params). */
export function mtdRangeUTC(now = new Date()): { from: string; to: string } {
  const y = now.getUTCFullYear();
  const m = now.getUTCMonth();
  const from = new Date(Date.UTC(y, m, 1));
  const to = new Date(Date.UTC(y, m, now.getUTCDate()));
  const fmt = (d: Date) => d.toISOString().slice(0, 10);
  return { from: fmt(from), to: fmt(to) };
}

export function currentPeriodUTC(now = new Date()): string {
  const y = now.getUTCFullYear();
  const m = String(now.getUTCMonth() + 1).padStart(2, "0");
  return `${y}-${m}`;
}

export function parseUsd(raw: string | undefined | null): number {
  if (!raw) return 0;
  const n = Number.parseFloat(raw);
  return Number.isFinite(n) ? n : 0;
}

export function usagePct(spend: number, budget: number): number {
  if (budget <= 0) return 0;
  return (spend / budget) * 100;
}

export function softThresholdBreached(
  spend: number,
  budget: number,
  thresholdPct: number,
): boolean {
  if (budget <= 0 || thresholdPct <= 0) return false;
  return spend * 100 >= budget * thresholdPct;
}

export type TeamSpendRow = {
  id?: string;
  name?: string;
  cost_usd?: string;
};

/** Match analytics team spend to a management team id/name. */
export function spendForTeam(
  teamId: string,
  teamName: string | undefined,
  rows: TeamSpendRow[],
): number {
  const byId = rows.find((r) => r.id && r.id === teamId);
  if (byId) return parseUsd(byId.cost_usd);
  if (teamName) {
    const byName = rows.find((r) => r.name && r.name.toLowerCase() === teamName.toLowerCase());
    if (byName) return parseUsd(byName.cost_usd);
  }
  return 0;
}
