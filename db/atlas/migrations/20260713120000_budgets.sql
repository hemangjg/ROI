-- Phase 2: team monthly budgets + soft-limit alerts
CREATE TABLE "public"."budgets" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "org_id" uuid NOT NULL,
  "team_id" uuid NOT NULL,
  "amount_usd" numeric(20,8) NOT NULL,
  "soft_threshold_pct" integer NOT NULL DEFAULT 80,
  "period_start" date NOT NULL,
  "created_by" uuid NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT "budgets_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "budgets_amount_positive" CHECK (amount_usd > 0),
  CONSTRAINT "budgets_threshold_range" CHECK ((soft_threshold_pct >= 1) AND (soft_threshold_pct <= 100)),
  CONSTRAINT "budgets_team_period_unique" UNIQUE ("team_id", "period_start"),
  CONSTRAINT "budgets_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organizations" ("id") ON DELETE CASCADE,
  CONSTRAINT "budgets_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "public"."teams" ("id") ON DELETE CASCADE,
  CONSTRAINT "budgets_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."users" ("id") ON DELETE SET NULL
);

CREATE INDEX "budgets_org_period_idx" ON "public"."budgets" ("org_id", "period_start" DESC);

CREATE TABLE "public"."budget_alerts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "budget_id" uuid NOT NULL,
  "org_id" uuid NOT NULL,
  "team_id" uuid NOT NULL,
  "threshold_pct" integer NOT NULL,
  "spend_usd" numeric(20,8) NOT NULL,
  "budget_usd" numeric(20,8) NOT NULL,
  "status" text NOT NULL DEFAULT 'open',
  "message" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT "budget_alerts_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "budget_alerts_status_check" CHECK (status IN ('open', 'acknowledged')),
  CONSTRAINT "budget_alerts_threshold_range" CHECK ((threshold_pct >= 1) AND (threshold_pct <= 100)),
  CONSTRAINT "budget_alerts_budget_threshold_unique" UNIQUE ("budget_id", "threshold_pct"),
  CONSTRAINT "budget_alerts_budget_id_fkey" FOREIGN KEY ("budget_id") REFERENCES "public"."budgets" ("id") ON DELETE CASCADE,
  CONSTRAINT "budget_alerts_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organizations" ("id") ON DELETE CASCADE,
  CONSTRAINT "budget_alerts_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "public"."teams" ("id") ON DELETE CASCADE
);

CREATE INDEX "budget_alerts_org_created_idx" ON "public"."budget_alerts" ("org_id", "created_at" DESC);
CREATE INDEX "budget_alerts_org_status_idx" ON "public"."budget_alerts" ("org_id", "status");
