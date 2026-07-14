-- name: CreateOrganization :one
INSERT INTO organizations (name, slug, timezone)
VALUES ($1, $2, COALESCE($3, 'UTC'))
RETURNING id, name, slug, timezone, created_at, updated_at;

-- name: GetOrganizationByID :one
SELECT id, name, slug, timezone, created_at, updated_at
FROM organizations
WHERE id = $1;

-- name: GetOrganizationBySlug :one
SELECT id, name, slug, timezone, created_at, updated_at
FROM organizations
WHERE slug = $1;

-- name: CreateOrgMembership :one
INSERT INTO org_memberships (user_id, org_id, role)
VALUES ($1, $2, $3)
RETURNING id, user_id, org_id, role, created_at;

-- name: ListOrgMembershipsByUser :many
SELECT id, user_id, org_id, role, created_at
FROM org_memberships
WHERE user_id = $1
ORDER BY created_at;

-- name: CreateTeam :one
INSERT INTO teams (org_id, name)
VALUES ($1, $2)
RETURNING id, org_id, name, deleted_at, created_at, updated_at;

-- name: ListTeamsByOrg :many
SELECT id, org_id, name, deleted_at, created_at, updated_at
FROM teams
WHERE org_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: SoftDeleteTeam :exec
UPDATE teams
SET deleted_at = now(), updated_at = now()
WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL;

-- name: CreateApiKey :one
INSERT INTO api_keys (org_id, name, key_prefix, key_hash, scopes, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, org_id, name, key_prefix, key_hash, scopes, created_by, revoked_at, created_at;

-- name: ListApiKeysByOrg :many
SELECT id, org_id, name, key_prefix, key_hash, scopes, created_by, revoked_at, created_at
FROM api_keys
WHERE org_id = $1 AND revoked_at IS NULL
ORDER BY created_at DESC;

-- name: RevokeApiKey :exec
UPDATE api_keys
SET revoked_at = now()
WHERE id = $1 AND org_id = $2 AND revoked_at IS NULL;

-- name: InsertAuditLog :one
INSERT INTO audit_logs (org_id, actor_id, actor_type, action, resource_type, resource_id, metadata, ip_address)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, org_id, actor_id, actor_type, action, resource_type, resource_id, metadata, ip_address, created_at;

-- name: ListAuditLogsByOrg :many
SELECT id, org_id, actor_id, actor_type, action, resource_type, resource_id, metadata, ip_address, created_at
FROM audit_logs
WHERE org_id = $1
  AND (
    sqlc.narg('cursor_created_at')::timestamptz IS NULL
    OR created_at < sqlc.narg('cursor_created_at')::timestamptz
  )
ORDER BY created_at DESC
LIMIT sqlc.arg('page_limit');

-- name: UpdateOrganization :one
UPDATE organizations
SET
  name = COALESCE(sqlc.narg('name'), name),
  timezone = COALESCE(sqlc.narg('timezone'), timezone),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING id, name, slug, timezone, created_at, updated_at;

-- name: GetOrgMembership :one
SELECT id, user_id, org_id, role, created_at
FROM org_memberships
WHERE user_id = $1 AND org_id = $2;

-- name: ListOrgMembersByOrg :many
SELECT
  u.id,
  u.email,
  u.name,
  m.role,
  m.created_at AS joined_at
FROM org_memberships m
JOIN users u ON u.id = m.user_id
WHERE m.org_id = $1
ORDER BY m.created_at;

-- name: GetTeamByID :one
SELECT id, org_id, name, deleted_at, created_at, updated_at
FROM teams
WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL;

-- name: GetApiKeyByID :one
SELECT id, org_id, name, key_prefix, key_hash, scopes, created_by, revoked_at, created_at
FROM api_keys
WHERE id = $1 AND org_id = $2 AND revoked_at IS NULL;

-- name: UpsertBudget :one
INSERT INTO budgets (org_id, team_id, amount_usd, soft_threshold_pct, period_start, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (team_id, period_start) DO UPDATE
SET
  amount_usd = EXCLUDED.amount_usd,
  soft_threshold_pct = EXCLUDED.soft_threshold_pct,
  updated_at = now()
RETURNING id, org_id, team_id, amount_usd, soft_threshold_pct, period_start, created_by, created_at, updated_at;

-- name: GetBudgetByTeamPeriod :one
SELECT id, org_id, team_id, amount_usd, soft_threshold_pct, period_start, created_by, created_at, updated_at
FROM budgets
WHERE org_id = $1 AND team_id = $2 AND period_start = $3;

-- name: ListBudgetsByOrgPeriod :many
SELECT b.id, b.org_id, b.team_id, t.name AS team_name, b.amount_usd, b.soft_threshold_pct,
       b.period_start, b.created_by, b.created_at, b.updated_at
FROM budgets b
JOIN teams t ON t.id = b.team_id
WHERE b.org_id = $1 AND b.period_start = $2
ORDER BY t.name;

-- name: InsertBudgetAlert :one
INSERT INTO budget_alerts (budget_id, org_id, team_id, threshold_pct, spend_usd, budget_usd, status, message)
VALUES ($1, $2, $3, $4, $5, $6, 'open', $7)
ON CONFLICT (budget_id, threshold_pct) DO NOTHING
RETURNING id, budget_id, org_id, team_id, threshold_pct, spend_usd, budget_usd, status, message, created_at;

-- name: ListBudgetAlertsByOrg :many
SELECT id, budget_id, org_id, team_id, threshold_pct, spend_usd, budget_usd, status, message, created_at
FROM budget_alerts
WHERE org_id = $1
  AND (sqlc.narg('status_filter')::text IS NULL OR status = sqlc.narg('status_filter'))
ORDER BY created_at DESC
LIMIT sqlc.arg('page_limit');

-- name: AcknowledgeBudgetAlert :one
UPDATE budget_alerts
SET status = 'acknowledged'
WHERE id = $1 AND org_id = $2 AND status = 'open'
RETURNING id, budget_id, org_id, team_id, threshold_pct, spend_usd, budget_usd, status, message, created_at;