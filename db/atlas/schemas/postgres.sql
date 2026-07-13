-- AI FinOps PostgreSQL declarative schema (Atlas source of truth)
-- Phase 1b (Ingestion MVP): business entities for auth, management, pricing, ingestion, workflow.

CREATE TABLE organizations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL,
    timezone    TEXT NOT NULL DEFAULT 'UTC',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT organizations_slug_unique UNIQUE (slug)
);

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    name          TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE TABLE org_memberships (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    org_id     UUID NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    role       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT org_memberships_role_check CHECK (role IN ('admin', 'member', 'viewer')),
    CONSTRAINT org_memberships_user_org_unique UNIQUE (user_id, org_id)
);

CREATE INDEX org_memberships_user_org_idx ON org_memberships (user_id, org_id);

CREATE TABLE teams (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id     UUID NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX teams_org_active_idx ON teams (org_id) WHERE deleted_at IS NULL;

CREATE TABLE team_members (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id    UUID NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT team_members_role_check CHECK (role IN ('manager', 'member')),
    CONSTRAINT team_members_team_user_unique UNIQUE (team_id, user_id)
);

CREATE TABLE api_keys (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      UUID NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    key_prefix  TEXT NOT NULL,
    key_hash    TEXT NOT NULL,
    scopes      TEXT[] NOT NULL DEFAULT ARRAY['ingest']::TEXT[],
    created_by  UUID REFERENCES users (id) ON DELETE SET NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX api_keys_prefix_active_idx ON api_keys (key_prefix) WHERE revoked_at IS NULL;

CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT refresh_tokens_token_hash_unique UNIQUE (token_hash)
);

CREATE INDEX refresh_tokens_user_id_idx ON refresh_tokens (user_id);

CREATE TABLE providers (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug       TEXT NOT NULL,
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT providers_slug_unique UNIQUE (slug)
);

CREATE TABLE models (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES providers (id) ON DELETE CASCADE,
    slug        TEXT NOT NULL,
    name        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT models_provider_slug_unique UNIQUE (provider_id, slug)
);

CREATE TABLE model_pricing (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id                 UUID NOT NULL REFERENCES models (id) ON DELETE CASCADE,
    input_price_per_1m       NUMERIC(20, 8) NOT NULL,
    output_price_per_1m      NUMERIC(20, 8) NOT NULL,
    cache_read_price_per_1m  NUMERIC(20, 8) NOT NULL DEFAULT 0,
    cache_write_price_per_1m NUMERIC(20, 8) NOT NULL DEFAULT 0,
    effective_from           TIMESTAMPTZ NOT NULL,
    effective_to             TIMESTAMPTZ,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX model_pricing_model_effective_idx ON model_pricing (model_id, effective_from DESC);

CREATE TABLE outbox_events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id       UUID NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    event_type   TEXT NOT NULL,
    payload      JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);

CREATE INDEX outbox_events_unpublished_idx ON outbox_events (created_at) WHERE published_at IS NULL;

CREATE TABLE idempotency_keys (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    idempotency_key TEXT NOT NULL,
    event_id        UUID NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL,
    CONSTRAINT idempotency_keys_org_key_unique UNIQUE (org_id, idempotency_key)
);

CREATE TABLE audit_logs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id        UUID NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    actor_id      UUID REFERENCES users (id) ON DELETE SET NULL,
    actor_type    TEXT NOT NULL DEFAULT 'user',
    action        TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id   TEXT,
    metadata      JSONB NOT NULL DEFAULT '{}'::JSONB,
    ip_address    INET,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT audit_logs_actor_type_check CHECK (actor_type IN ('user', 'system', 'api_key'))
);

CREATE INDEX audit_logs_org_created_idx ON audit_logs (org_id, created_at DESC);

CREATE OR REPLACE FUNCTION prevent_audit_log_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'audit_logs are append-only';
END;
$$;

CREATE TRIGGER audit_logs_no_update
    BEFORE UPDATE ON audit_logs
    FOR EACH ROW
    EXECUTE FUNCTION prevent_audit_log_mutation();

CREATE TRIGGER audit_logs_no_delete
    BEFORE DELETE ON audit_logs
    FOR EACH ROW
    EXECUTE FUNCTION prevent_audit_log_mutation();