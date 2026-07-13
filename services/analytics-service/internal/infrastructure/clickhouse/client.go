package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DailyPoint is a single day's aggregated spend.
type DailyPoint struct {
	Date    time.Time
	CostUSD decimal.Decimal
}

// ProviderSpend is spend grouped by provider.
type ProviderSpend struct {
	Provider   string
	CostUSD    decimal.Decimal
	EventCount uint64
}

// TeamSpend is spend grouped by team.
type TeamSpend struct {
	TeamID     *uuid.UUID
	CostUSD    decimal.Decimal
	EventCount uint64
}

// Client reads analytics rollups from ClickHouse.
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

// SumCost returns total spend for an org in the inclusive date range.
func (c *Client) SumCost(ctx context.Context, orgID uuid.UUID, from, to time.Time, teamID *uuid.UUID) (decimal.Decimal, error) {
	query := `
		SELECT coalesce(sum(total_cost_usd), 0)
		FROM finops.daily_spend_rollups
		WHERE org_id = ?
		  AND rollup_date >= ?
		  AND rollup_date <= ?
	`
	args := []any{orgID, from.UTC(), to.UTC()}
	if teamID != nil {
		query += " AND team_id = ?"
		args = append(args, *teamID)
	}

	var total decimal.Decimal
	if err := c.conn.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return decimal.Zero, fmt.Errorf("sum cost: %w", err)
	}
	return total, nil
}

// DailyCosts returns per-day spend ordered ascending.
func (c *Client) DailyCosts(ctx context.Context, orgID uuid.UUID, from, to time.Time, teamID *uuid.UUID) ([]DailyPoint, error) {
	query := `
		SELECT rollup_date, sum(total_cost_usd) AS cost_usd
		FROM finops.daily_spend_rollups
		WHERE org_id = ?
		  AND rollup_date >= ?
		  AND rollup_date <= ?
	`
	args := []any{orgID, from.UTC(), to.UTC()}
	if teamID != nil {
		query += " AND team_id = ?"
		args = append(args, *teamID)
	}
	query += " GROUP BY rollup_date ORDER BY rollup_date ASC"

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("daily costs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	points := make([]DailyPoint, 0)
	for rows.Next() {
		var point DailyPoint
		if err := rows.Scan(&point.Date, &point.CostUSD); err != nil {
			return nil, fmt.Errorf("scan daily cost: %w", err)
		}
		points = append(points, point)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate daily costs: %w", err)
	}
	return points, nil
}

// SpendByProvider returns provider breakdown for the date range.
func (c *Client) SpendByProvider(ctx context.Context, orgID uuid.UUID, from, to time.Time, teamID *uuid.UUID) ([]ProviderSpend, error) {
	query := `
		SELECT provider, sum(total_cost_usd) AS cost_usd, sum(event_count) AS event_count
		FROM finops.daily_spend_rollups
		WHERE org_id = ?
		  AND rollup_date >= ?
		  AND rollup_date <= ?
	`
	args := []any{orgID, from.UTC(), to.UTC()}
	if teamID != nil {
		query += " AND team_id = ?"
		args = append(args, *teamID)
	}
	query += " GROUP BY provider ORDER BY cost_usd DESC"

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("spend by provider: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]ProviderSpend, 0)
	for rows.Next() {
		var item ProviderSpend
		if err := rows.Scan(&item.Provider, &item.CostUSD, &item.EventCount); err != nil {
			return nil, fmt.Errorf("scan provider spend: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate provider spend: %w", err)
	}
	return items, nil
}

// SpendByTeam returns team breakdown for the date range.
func (c *Client) SpendByTeam(ctx context.Context, orgID uuid.UUID, from, to time.Time) ([]TeamSpend, error) {
	rows, err := c.conn.Query(ctx, `
		SELECT team_id, sum(total_cost_usd) AS cost_usd, sum(event_count) AS event_count
		FROM finops.daily_spend_rollups
		WHERE org_id = ?
		  AND rollup_date >= ?
		  AND rollup_date <= ?
		GROUP BY team_id
		ORDER BY cost_usd DESC
	`, orgID, from.UTC(), to.UTC())
	if err != nil {
		return nil, fmt.Errorf("spend by team: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]TeamSpend, 0)
	for rows.Next() {
		var item TeamSpend
		if err := rows.Scan(&item.TeamID, &item.CostUSD, &item.EventCount); err != nil {
			return nil, fmt.Errorf("scan team spend: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate team spend: %w", err)
	}
	return items, nil
}

// Close closes the ClickHouse connection.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}