-- name: GetProviderBySlug :one
SELECT id, slug, name, created_at
FROM providers
WHERE slug = $1;

-- name: GetModelByProviderAndSlug :one
SELECT m.id, m.provider_id, m.slug, m.name, m.created_at
FROM models m
JOIN providers p ON p.id = m.provider_id
WHERE p.slug = $1 AND m.slug = $2;

-- name: GetPriceAtTime :one
SELECT
    mp.id,
    mp.model_id,
    mp.input_price_per_1m,
    mp.output_price_per_1m,
    mp.cache_read_price_per_1m,
    mp.cache_write_price_per_1m,
    mp.effective_from,
    mp.effective_to,
    mp.created_at
FROM model_pricing mp
JOIN models m ON m.id = mp.model_id
JOIN providers p ON p.id = m.provider_id
WHERE p.slug = $1
  AND m.slug = $2
  AND mp.effective_from <= $3
  AND (mp.effective_to IS NULL OR mp.effective_to > $3)
ORDER BY mp.effective_from DESC
LIMIT 1;

-- name: UpsertProvider :one
INSERT INTO providers (slug, name)
VALUES ($1, $2)
ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
RETURNING id, slug, name, created_at;

-- name: UpsertModel :one
INSERT INTO models (provider_id, slug, name)
VALUES ($1, $2, $3)
ON CONFLICT (provider_id, slug) DO UPDATE SET name = EXCLUDED.name
RETURNING id, provider_id, slug, name, created_at;

-- name: InsertModelPricing :one
INSERT INTO model_pricing (
    model_id,
    input_price_per_1m,
    output_price_per_1m,
    cache_read_price_per_1m,
    cache_write_price_per_1m,
    effective_from,
    effective_to
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, model_id, input_price_per_1m, output_price_per_1m, cache_read_price_per_1m, cache_write_price_per_1m, effective_from, effective_to, created_at;

-- name: CloseOpenModelPricing :exec
UPDATE model_pricing
SET effective_to = $2
WHERE model_id = $1 AND effective_to IS NULL;