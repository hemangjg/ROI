-- Phase 1b (Ingestion MVP): business schema for auth, management, pricing, ingestion, workflow.

CREATE TABLE "public"."organizations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  "slug" text NOT NULL,
  "timezone" text NOT NULL DEFAULT 'UTC',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "organizations_slug_unique" UNIQUE ("slug")
);

CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "email" text NOT NULL,
  "password_hash" text NOT NULL,
  "name" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_unique" UNIQUE ("email")
);

CREATE TABLE "public"."org_memberships" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "org_id" uuid NOT NULL,
  "role" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "org_memberships_role_check" CHECK (role IN ('admin', 'member', 'viewer')),
  CONSTRAINT "org_memberships_user_org_unique" UNIQUE ("user_id", "org_id"),
  CONSTRAINT "org_memberships_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "org_memberships_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "org_memberships_user_org_idx" ON "public"."org_memberships" ("user_id", "org_id");

CREATE TABLE "public"."teams" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "org_id" uuid NOT NULL,
  "name" text NOT NULL,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "teams_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "teams_org_active_idx" ON "public"."teams" ("org_id") WHERE (deleted_at IS NULL);

CREATE TABLE "public"."team_members" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "team_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "role" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "team_members_role_check" CHECK (role IN ('manager', 'member')),
  CONSTRAINT "team_members_team_user_unique" UNIQUE ("team_id", "user_id"),
  CONSTRAINT "team_members_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "public"."teams" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "team_members_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE "public"."api_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "org_id" uuid NOT NULL,
  "name" text NOT NULL,
  "key_prefix" text NOT NULL,
  "key_hash" text NOT NULL,
  "scopes" text[] NOT NULL DEFAULT ARRAY['ingest']::text[],
  "created_by" uuid NULL,
  "revoked_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "api_keys_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "api_keys_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);

CREATE INDEX "api_keys_prefix_active_idx" ON "public"."api_keys" ("key_prefix") WHERE (revoked_at IS NULL);

CREATE TABLE "public"."refresh_tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "token_hash" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "revoked_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "refresh_tokens_token_hash_unique" UNIQUE ("token_hash"),
  CONSTRAINT "refresh_tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "refresh_tokens_user_id_idx" ON "public"."refresh_tokens" ("user_id");

CREATE TABLE "public"."providers" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "slug" text NOT NULL,
  "name" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "providers_slug_unique" UNIQUE ("slug")
);

CREATE TABLE "public"."models" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "provider_id" uuid NOT NULL,
  "slug" text NOT NULL,
  "name" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "models_provider_slug_unique" UNIQUE ("provider_id", "slug"),
  CONSTRAINT "models_provider_id_fkey" FOREIGN KEY ("provider_id") REFERENCES "public"."providers" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE "public"."model_pricing" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "model_id" uuid NOT NULL,
  "input_price_per_1m" numeric(20,8) NOT NULL,
  "output_price_per_1m" numeric(20,8) NOT NULL,
  "cache_read_price_per_1m" numeric(20,8) NOT NULL DEFAULT 0,
  "cache_write_price_per_1m" numeric(20,8) NOT NULL DEFAULT 0,
  "effective_from" timestamptz NOT NULL,
  "effective_to" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "model_pricing_model_id_fkey" FOREIGN KEY ("model_id") REFERENCES "public"."models" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "model_pricing_model_effective_idx" ON "public"."model_pricing" ("model_id", "effective_from" DESC);

CREATE TABLE "public"."outbox_events" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "org_id" uuid NOT NULL,
  "event_type" text NOT NULL,
  "payload" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "published_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "outbox_events_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "outbox_events_unpublished_idx" ON "public"."outbox_events" ("created_at") WHERE (published_at IS NULL);

CREATE TABLE "public"."idempotency_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "org_id" uuid NOT NULL,
  "idempotency_key" text NOT NULL,
  "event_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "expires_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "idempotency_keys_org_key_unique" UNIQUE ("org_id", "idempotency_key"),
  CONSTRAINT "idempotency_keys_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE "public"."audit_logs" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "org_id" uuid NOT NULL,
  "actor_id" uuid NULL,
  "actor_type" text NOT NULL DEFAULT 'user',
  "action" text NOT NULL,
  "resource_type" text NOT NULL,
  "resource_id" text NULL,
  "metadata" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "ip_address" inet NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "audit_logs_actor_type_check" CHECK (actor_type IN ('user', 'system', 'api_key')),
  CONSTRAINT "audit_logs_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "audit_logs_actor_id_fkey" FOREIGN KEY ("actor_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);

CREATE INDEX "audit_logs_org_created_idx" ON "public"."audit_logs" ("org_id", "created_at" DESC);

CREATE OR REPLACE FUNCTION "public"."prevent_audit_log_mutation"() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    RAISE EXCEPTION 'audit_logs are append-only';
END;
$$;

CREATE TRIGGER "audit_logs_no_update" BEFORE UPDATE ON "public"."audit_logs" FOR EACH ROW EXECUTE FUNCTION "public"."prevent_audit_log_mutation"();
CREATE TRIGGER "audit_logs_no_delete" BEFORE DELETE ON "public"."audit_logs" FOR EACH ROW EXECUTE FUNCTION "public"."prevent_audit_log_mutation"();