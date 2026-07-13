package clickhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
)

// UsageEventRow is a ClickHouse usage_events insert payload.
type UsageEventRow struct {
	ID               uuid.UUID
	OrgID            uuid.UUID
	IdempotencyKey   string
	TeamID           *uuid.UUID
	UserID           *uuid.UUID
	Provider         string
	Model            string
	InputTokens      uint32
	OutputTokens     uint32
	CacheReadTokens  uint32
	CacheWriteTokens uint32
	CostUSD          string
	IsPriced         bool
	Project          string
	Feature          string
	Customer         string
	Metadata         map[string]string
	OccurredAt       time.Time
	ReceivedAt       time.Time
}

// Client inserts workflow-processed usage events.
type Client struct {
	conn clickhouse.Conn
}

// NewClient opens a ClickHouse connection.
func NewClient(url string) (*Client, error) {
	opts, err := clickhouse.ParseDSN(url)
	if err != nil {
		return nil, fmt.Errorf("parse clickhouse dsn: %w", err)
	}
	conn, err := clickhouse.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("open clickhouse: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Ping(ctx); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping clickhouse: %w", err)
	}
	return &Client{conn: conn}, nil
}

// InsertUsageEvent writes a processed event to finops.usage_events.
func (c *Client) InsertUsageEvent(ctx context.Context, row UsageEventRow) error {
	metadataJSON, err := json.Marshal(row.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	isPriced := uint8(0)
	if row.IsPriced {
		isPriced = 1
	}

	// Nullable UUID columns must not receive a typed nil *uuid.UUID — the driver
	// panics calling Value on a nil pointer. Use untyped nil / concrete values.
	var teamID any
	if row.TeamID != nil {
		teamID = *row.TeamID
	}
	var userID any
	if row.UserID != nil {
		userID = *row.UserID
	}

	err = c.conn.Exec(ctx, `
		INSERT INTO finops.usage_events (
			id, org_id, idempotency_key, team_id, user_id, provider, model,
			input_tokens, output_tokens, cache_read_tokens, cache_write_tokens,
			cost_usd, is_priced, project, feature, customer, metadata,
			occurred_at, received_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		row.ID,
		row.OrgID,
		row.IdempotencyKey,
		teamID,
		userID,
		row.Provider,
		row.Model,
		row.InputTokens,
		row.OutputTokens,
		row.CacheReadTokens,
		row.CacheWriteTokens,
		row.CostUSD,
		isPriced,
		row.Project,
		row.Feature,
		row.Customer,
		string(metadataJSON),
		row.OccurredAt.UTC(),
		row.ReceivedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert usage event: %w", err)
	}
	return nil
}

// Close closes the ClickHouse connection.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}