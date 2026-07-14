"use client";

import { useEffect, useMemo, useState } from "react";
import { Wallet } from "lucide-react";

import { AppHeader } from "@/components/layout/app-header";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api, type Budget, type BudgetAlert, type Team } from "@/lib/api/client";
import { useOrgId } from "@/hooks/use-org-id";

function currentPeriod(): string {
  const d = new Date();
  return `${d.getUTCFullYear()}-${String(d.getUTCMonth() + 1).padStart(2, "0")}`;
}

export default function BudgetsPage() {
  const orgId = useOrgId();
  const period = useMemo(() => currentPeriod(), []);
  const [teams, setTeams] = useState<Team[]>([]);
  const [budgets, setBudgets] = useState<Budget[]>([]);
  const [alerts, setAlerts] = useState<BudgetAlert[]>([]);
  const [teamId, setTeamId] = useState("");
  const [amount, setAmount] = useState("100.00");
  const [threshold, setThreshold] = useState("80");
  const [spend, setSpend] = useState("85.00");
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  async function refresh() {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const [t, b, a] = await Promise.all([
        api.management.listTeams(orgId),
        api.management.listBudgets(orgId, period),
        api.management.listBudgetAlerts(orgId, "open"),
      ]);
      setTeams(t);
      setBudgets(b);
      setAlerts(a);
      if (!teamId && t.length > 0) setTeamId(t[0].id);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load budgets");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [orgId, period]);

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

  async function onEvaluate() {
    if (!orgId || !teamId) return;
    setMessage(null);
    setError(null);
    try {
      const status = await api.management.evaluateBudget(orgId, {
        team_id: teamId,
        spend_usd: spend,
        period_start: period,
      });
      setMessage(
        status.over_soft_threshold
          ? `Soft limit hit (${status.usage_pct.toFixed(1)}%). Alert ${status.alert_fired ? "created" : "already open"}.`
          : `Within budget (${status.usage_pct.toFixed(1)}% used).`,
      );
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to evaluate budget");
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
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Wallet className="h-5 w-5" />
              Team monthly budgets
            </CardTitle>
            <CardDescription>
              Period {period} (UTC). Soft alerts fire when spend reaches the threshold (default 80%).
            </CardDescription>
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
                <Input id="threshold" value={threshold} onChange={(e) => setThreshold(e.target.value)} />
              </div>
              <Button onClick={() => void onSaveBudget()} disabled={!teamId || loading}>
                Save budget
              </Button>
            </div>
            <div className="space-y-3">
              <div className="space-y-1">
                <Label htmlFor="spend">Reported MTD spend (USD)</Label>
                <Input id="spend" value={spend} onChange={(e) => setSpend(e.target.value)} />
              </div>
              <p className="text-sm text-muted-foreground">
                Evaluate against the saved budget. Wire this to live analytics spend in a follow-up.
              </p>
              <Button variant="secondary" onClick={() => void onEvaluate()} disabled={!teamId || loading}>
                Evaluate soft limit
              </Button>
              {message ? <p className="text-sm text-foreground">{message}</p> : null}
              {error ? <p className="text-sm text-destructive">{error}</p> : null}
            </div>
          </CardContent>
        </Card>

        <div className="grid gap-6 lg:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Configured budgets</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              {loading ? <p className="text-muted-foreground">Loading…</p> : null}
              {!loading && budgets.length === 0 ? (
                <p className="text-muted-foreground">No budgets for this period yet.</p>
              ) : null}
              {budgets.map((b) => (
                <div key={b.id} className="flex items-center justify-between rounded-md border px-3 py-2">
                  <div>
                    <p className="font-medium">{b.team_name || b.team_id}</p>
                    <p className="text-muted-foreground">
                      ${b.amount_usd} · soft {b.soft_threshold_pct}%
                    </p>
                  </div>
                  <span className="text-xs text-muted-foreground">{b.period_start}</span>
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
        </div>
      </main>
    </>
  );
}
