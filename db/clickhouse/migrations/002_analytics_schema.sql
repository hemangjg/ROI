-- Phase 1b (Ingestion MVP): usage telemetry and dashboard rollups.

CREATE TABLE IF NOT EXISTS finops.usage_events (
    id UUID,
    org_id UUID,
    idempotency_key String,
    team_id Nullable(UUID),
    user_id Nullable(UUID),
    provider LowCardinality(String),
    model LowCardinality(String),
    input_tokens UInt32,
    output_tokens UInt32,
    cache_read_tokens UInt32 DEFAULT 0,
    cache_write_tokens UInt32 DEFAULT 0,
    cost_usd Decimal64(8),
    is_priced UInt8,
    project String DEFAULT '',
    feature String DEFAULT '',
    customer String DEFAULT '',
    metadata String,
    occurred_at DateTime64(3, 'UTC'),
    received_at DateTime64(3, 'UTC')
) ENGINE = MergeTree()
PARTITION BY (org_id, toYYYYMM(occurred_at))
ORDER BY (org_id, occurred_at, id)
TTL toDate(occurred_at) + INTERVAL 90 DAY;

CREATE TABLE IF NOT EXISTS finops.daily_spend_rollups (
    org_id UUID,
    rollup_date Date,
    team_id Nullable(UUID),
    provider LowCardinality(String),
    model LowCardinality(String),
    total_input_tokens UInt64,
    total_output_tokens UInt64,
    total_cost_usd Decimal64(8),
    event_count UInt32
) ENGINE = SummingMergeTree()
PARTITION BY (org_id, toYYYYMM(rollup_date))
ORDER BY (org_id, rollup_date, team_id, provider, model)
SETTINGS allow_nullable_key = 1;

CREATE MATERIALIZED VIEW IF NOT EXISTS finops.daily_spend_rollups_mv
TO finops.daily_spend_rollups AS
SELECT
    org_id,
    toDate(occurred_at) AS rollup_date,
    team_id,
    provider,
    model,
    sum(input_tokens) AS total_input_tokens,
    sum(output_tokens) AS total_output_tokens,
    sum(cost_usd) AS total_cost_usd,
    count() AS event_count
FROM finops.usage_events
GROUP BY org_id, rollup_date, team_id, provider, model;