"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { AlertTriangle, RefreshCw, Wallet } from "lucide-react";

import { AppHeader } from "@/components/layout/app-header";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useOrgId } from "@/hooks/use-org-id";
import { api, type Budget, type BudgetAlert, type Team } from "@/lib/api/client";
import {
  currentPeriodUTC,
  mtdRangeUTC,
  parseUsd,
  softThresholdBreached,
  spendForTeam,
  usagePct,
} from "@/lib/budgets";
import { cn } from "@/lib/utils";

type RowView = Budget & {
  spend: number;
  pct: number;
  overSoft: boolean;
};

function BudgetBar({ spend, budget, softPct }: { spend: number; budget: number; softPct: number }) {
  const used = budget > 0 ? Math.min(100, (spend / budget) * 100) : 0;
  const softLeft = Math.min(100, Math.max(0, softPct));
  const over = used >= softLeft && softLeft > 0;

  return (
    <div className="space-y-1">
      <div className="relative h-2.5 w-full overflow-hidden rounded-full bg-muted">
        <div
          className={
            over
              ? "absolute inset-y-0 left-0 rounded-full bg-destructive transition-all"
              : "absolute inset-y-0 left-0 rounded-full bg-primary transition-all"
          }
          style={{ width: `${used}%` }}
        />
        <div
          className="absolute inset-y-0 w-0.5 bg-amber-500"
          style={{ left: `${softLeft}%` }}
          title={`Soft ${softPct}%`}
        />
      </div>
      <div className="flex justify-between text-xs text-muted-foreground">
        <span>
          ${spend.toFixed(2)} / ${budget.toFixed(2)} ({used.toFixed(1)}%)
        </span>
        <span className="text-amber-600">soft {softPct}%</span>
      </div>
    </div>
  );
}

export default function BudgetsPage() {
  const orgId = useOrgId();
  const period = useMemo(() => currentPeriodUTC(), []);
  const mtd = useMemo(() => mtdRangeUTC(), []);

  const [teams, setTeams] = useState<Team[]>([]);
  const [budgets, setBudgets] = useState<Budget[]>([]);
  const [alerts, setAlerts] = useState<BudgetAlert[]>([]);
  const [rows, setRows] = useState<RowView[]>([]);
  const [teamId, setTeamId] = useState("");
  const [amount, setAmount] = useState("100.00");
  const [threshold, setThreshold] = useState("80");
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [evaluating, setEvaluating] = useState(false);

  const refresh = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const [t, b, a, spend] = await Promise.all([
        api.management.listTeams(orgId),
        api.management.listBudgets(orgId, period),
        api.management.listBudgetAlerts(orgId, "open"),
        api.analytics.getSpendByTeam(orgId, mtd),
      ]);
      setTeams(t);
      setBudgets(b);
      setAlerts(a);
      const firstTeam = t.find((x) => Boolean(x.id));
      if (!teamId && firstTeam) setTeamId(firstTeam.id);

      const teamSpend = spend.teams ?? [];
      setRows(
        b.map((budget) => {
          const s = spendForTeam(budget.team_id, budget.team_name, teamSpend);
          const budgetAmt = parseUsd(budget.amount_usd);
          const pct = usagePct(s, budgetAmt);
          return {
            ...budget,
            spend: s,
            pct,
            overSoft: softThresholdBreached(s, budgetAmt, budget.soft_threshold_pct),
          };
        }),
      );
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load budgets");
    } finally {
      setLoading(false);
    }
  }, [orgId, period, mtd, teamId]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  async function onSaveBudget() {
    if (!orgId || !teamId) return;
    setMessage(null);
    setError(null);
    try {
      await api.management.upsertBudget(orgId, {
        team_id: teamId,
        amount_usd: amount,
        soft_threshold_pct: Number(threshold) || 80,
        period_start: period,
      });
      setMessage("Budget saved");
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to save budget");
    }
  }

  async function onEvaluateLive() {
    if (!orgId) return;
    setEvaluating(true);
    setMessage(null);
    setError(null);
    try {
      const spend = await api.analytics.getSpendByTeam(orgId, mtd);
      const teamSpend = spend.teams ?? [];
      const list = budgets.length ? budgets : await api.management.listBudgets(orgId, period);
      let fired = 0;
      let checked = 0;
      for (const b of list) {
        const s = spendForTeam(b.team_id, b.team_name, teamSpend);
        checked += 1;
        const status = await api.management.evaluateBudget(orgId, {
          team_id: b.team_id,
          spend_usd: s.toFixed(2),
          period_start: period,
        });
        if (status.alert_fired) fired += 1;
      }
      setMessage(
        checked === 0
          ? "No budgets to evaluate for this period."
          : `Evaluated ${checked} budget(s) against live MTD spend. New alerts: ${fired}.`,
      );
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to evaluate budgets");
    } finally {
      setEvaluating(false);
    }
  }

  async function onAck(alertId: string) {
    if (!orgId) return;
    try {
      await api.management.acknowledgeBudgetAlert(orgId, alertId);
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to acknowledge alert");
    }
  }

  return (
    <>
      <AppHeader title="Budgets" showDateRange={false} />
      <main className="space-y-6 px-8 py-6">
        <Card>
          <CardHeader className="flex flex-row items-start justify-between gap-4">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Wallet className="h-5 w-5" />
                Team monthly budgets
              </CardTitle>
              <CardDescription>
                Period {period} (UTC). Progress uses live MTD spend from analytics ({mtd.from} →{" "}
                {mtd.to}). Soft line marks alert threshold.
              </CardDescription>
            </div>
            <Button
              variant="secondary"
              size="sm"
              disabled={loading || evaluating}
              onClick={() => void onEvaluateLive()}
            >
              <RefreshCw className={cn("mr-2 h-4 w-4", evaluating && "animate-spin")} />
              Evaluate live spend
            </Button>
          </CardHeader>
          <CardContent className="grid gap-4 md:grid-cols-2">
            <div className="space-y-3">
              <div className="space-y-1">
                <Label htmlFor="team">Team</Label>
                <select
                  id="team"
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
                  value={teamId}
                  onChange={(e) => setTeamId(e.target.value)}
                >
                  {teams.map((t) => (
                    <option key={t.id} value={t.id}>
                      {t.name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-1">
                <Label htmlFor="amount">Budget (USD)</Label>
                <Input id="amount" value={amount} onChange={(e) => setAmount(e.target.value)} />
              </div>
              <div className="space-y-1">
                <Label htmlFor="threshold">Soft threshold %</Label>
                <Input
                  id="threshold"
                  value={threshold}
                  onChange={(e) => setThreshold(e.target.value)}
                />
              </div>
              <Button onClick={() => void onSaveBudget()} disabled={!teamId || loading}>
                Save budget
              </Button>
            </div>
            <div className="space-y-3 text-sm text-muted-foreground">
              <p>
                After saving, use <strong className="text-foreground">Evaluate live spend</strong>{" "}
                to check every budget against analytics MTD by-team totals and open soft alerts.
              </p>
              {message ? <p className="text-foreground">{message}</p> : null}
              {error ? <p className="text-destructive">{error}</p> : null}
              {alerts.length > 0 ? (
                <p className="flex items-center gap-2 text-amber-700">
                  <AlertTriangle className="h-4 w-4" />
                  {alerts.length} open alert{alerts.length === 1 ? "" : "s"}
                </p>
              ) : null}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Budget vs actual (MTD)</CardTitle>
            <CardDescription>
              Amber marker = soft threshold. Red fill = threshold reached or exceeded.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {loading ? <p className="text-sm text-muted-foreground">Loading…</p> : null}
            {!loading && rows.length === 0 ? (
              <p className="text-sm text-muted-foreground">No budgets for this period yet.</p>
            ) : null}
            {rows.map((r) => (
              <div key={r.id} className="space-y-2 rounded-md border px-4 py-3">
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <p className="font-medium">{r.team_name || r.team_id}</p>
                    <p className="text-xs text-muted-foreground">
                      budget ${parseUsd(r.amount_usd).toFixed(2)} · soft {r.soft_threshold_pct}%
                    </p>
                  </div>
                  {r.overSoft ? (
                    <span className="rounded-full bg-destructive/10 px-2 py-0.5 text-xs font-medium text-destructive">
                      Soft limit
                    </span>
                  ) : (
                    <span className="rounded-full bg-muted px-2 py-0.5 text-xs text-muted-foreground">
                      On track
                    </span>
                  )}
                </div>
                <BudgetBar
                  spend={r.spend}
                  budget={parseUsd(r.amount_usd)}
                  softPct={r.soft_threshold_pct}
                />
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Open alerts</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm">
            {alerts.length === 0 ? <p className="text-muted-foreground">No open alerts.</p> : null}
            {alerts.map((a) => (
              <div key={a.id} className="space-y-2 rounded-md border px-3 py-2">
                <p>{a.message}</p>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">
                    spend ${a.spend_usd} / ${a.budget_usd}
                  </span>
                  <Button size="sm" variant="outline" onClick={() => void onAck(a.id)}>
                    Acknowledge
                  </Button>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
      </main>
    </>
  );
}
