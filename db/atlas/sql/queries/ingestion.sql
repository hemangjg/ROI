-- name: GetApiKeyByPrefix :one
SELECT id, org_id, name, key_prefix, key_hash, scopes, created_by, revoked_at, created_at
FROM api_keys
WHERE key_prefix = $1 AND revoked_at IS NULL;

-- name: InsertOutboxEvent :one
INSERT INTO outbox_events (id, org_id, event_type, payload)
VALUES ($1, $2, $3, $4)
RETURNING id, org_id, event_type, payload, created_at, published_at;

-- name: ListUnpublishedOutboxEvents :many
SELECT id, org_id, event_type, payload, created_at, published_at
FROM outbox_events
WHERE published_at IS NULL
ORDER BY created_at
LIMIT $1;

-- name: MarkOutboxEventPublished :exec
UPDATE outbox_events
SET published_at = now()
WHERE id = $1 AND published_at IS NULL;

-- name: GetOutboxEventByID :one
SELECT id, org_id, event_type, payload, created_at, published_at
FROM outbox_events
WHERE id = $1 AND org_id = $2;