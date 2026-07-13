-- name: CheckIdempotency :one
SELECT id, org_id, idempotency_key, event_id, created_at, expires_at
FROM idempotency_keys
WHERE org_id = $1 AND idempotency_key = $2 AND expires_at > now();

-- name: RecordIdempotency :one
INSERT INTO idempotency_keys (org_id, idempotency_key, event_id, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id, org_id, idempotency_key, event_id, created_at, expires_at;