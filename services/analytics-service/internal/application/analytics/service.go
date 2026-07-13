package analytics

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	chinfra "github.com/ai-finops/ai-finops/services/analytics-service/internal/infrastructure/clickhouse"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

const defaultRangeDays = 30

// QueryParams captures analytics filter inputs.
type QueryParams struct {
	OrgID   uuid.UUID
	From    time.Time
	To      time.Time
	TeamID  *uuid.UUID
}

// SpendSummaryResult matches the OpenAPI SpendSummary schema.
type SpendSummaryResult struct {
	TodayUSD          string  `json:"today_usd"`
	MTDUSD            string  `json:"mtd_usd"`
	PreviousPeriodUSD string  `json:"previous_period_usd"`
	ChangePct         float64 `json:"change_pct"`
}

// TimeSeriesPoint is a single daily spend data point.
type TimeSeriesPoint struct {
	Date    string `json:"date"`
	CostUSD string `json:"cost_usd"`
}

// TimeSeriesResult matches the OpenAPI SpendTimeSeries schema.
type TimeSeriesResult struct {
	Data []TimeSeriesPoint `json:"data"`
}

// ProviderBreakdownItem is one provider row.
type ProviderBreakdownItem struct {
	Name    string  `json:"name"`
	CostUSD string  `json:"cost_usd"`
	Pct     float64 `json:"pct"`
}

// ProviderBreakdownResult matches the OpenAPI SpendByProvider schema.
type ProviderBreakdownResult struct {
	Providers []ProviderBreakdownItem `json:"providers"`
}

// TeamBreakdownItem is one team row.
type TeamBreakdownItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CostUSD    string `json:"cost_usd"`
	EventCount int    `json:"event_count"`
}

// TeamBreakdownResult matches the OpenAPI SpendByTeam schema.
type TeamBreakdownResult struct {
	Teams []TeamBreakdownItem `json:"teams"`
}

// Service implements analytics queries over ClickHouse rollups.
type Service struct {
	clickhouse *chinfra.Client
	queries    *sqlcgen.Queries
	log        *slog.Logger
}

// NewService constructs the analytics application service.
func NewService(clickhouse *chinfra.Client, queries *sqlcgen.Queries, log *slog.Logger) *Service {
	return &Service{
		clickhouse: clickhouse,
		queries:    queries,
		log:        log,
	}
}

// ParseQueryParams validates org and optional filters from HTTP/gRPC inputs.
func ParseQueryParams(orgID string, fromRaw, toRaw, teamIDRaw string) (QueryParams, error) {
	parsedOrgID, err := uuid.Parse(orgID)
	if err != nil {
		return QueryParams{}, shared.BadRequest("invalid_org_id", "org_id must be a valid uuid")
	}

	now := time.Now().UTC()
	to := endOfDay(now)
	from := startOfDay(now.AddDate(0, 0, -(defaultRangeDays - 1)))

	if fromRaw != "" {
		from, err = parseTimestamp(fromRaw)
		if err != nil {
			return QueryParams{}, shared.BadRequest("invalid_from", "from must be a valid ISO 8601 timestamp")
		}
		from = startOfDay(from)
	}
	if toRaw != "" {
		to, err = parseTimestamp(toRaw)
		if err != nil {
			return QueryParams{}, shared.BadRequest("invalid_to", "to must be a valid ISO 8601 timestamp")
		}
		to = endOfDay(to)
	}
	if to.Before(from) {
		return QueryParams{}, shared.BadRequest("invalid_range", "to must be on or after from")
	}

	var teamID *uuid.UUID
	if teamIDRaw != "" {
		parsedTeamID, err := uuid.Parse(teamIDRaw)
		if err != nil {
			return QueryParams{}, shared.BadRequest("invalid_team_id", "team_id must be a valid uuid")
		}
		teamID = &parsedTeamID
	}

	return QueryParams{
		OrgID:  parsedOrgID,
		From:   from,
		To:     to,
		TeamID: teamID,
	}, nil
}

// GetSpendSummary returns today, MTD, and period-over-period spend metrics.
func (s *Service) GetSpendSummary(ctx context.Context, params QueryParams) (SpendSummaryResult, error) {
	now := time.Now().UTC()
	today := startOfDay(now)

	todayUSD, err := s.clickhouse.SumCost(ctx, params.OrgID, today, today, params.TeamID)
	if err != nil {
		return SpendSummaryResult{}, shared.Internal("failed to query today spend")
	}

	mtdStart := startOfMonth(now)
	mtdUSD, err := s.clickhouse.SumCost(ctx, params.OrgID, mtdStart, today, params.TeamID)
	if err != nil {
		return SpendSummaryResult{}, shared.Internal("failed to query mtd spend")
	}

	currentUSD, err := s.clickhouse.SumCost(ctx, params.OrgID, params.From, params.To, params.TeamID)
	if err != nil {
		return SpendSummaryResult{}, shared.Internal("failed to query current period spend")
	}

	periodDays := int(params.To.Sub(params.From).Hours()/24) + 1
	prevTo := params.From.AddDate(0, 0, -1)
	prevFrom := prevTo.AddDate(0, 0, -(periodDays - 1))

	previousUSD, err := s.clickhouse.SumCost(ctx, params.OrgID, prevFrom, prevTo, params.TeamID)
	if err != nil {
		return SpendSummaryResult{}, shared.Internal("failed to query previous period spend")
	}

	return SpendSummaryResult{
		TodayUSD:          formatUSD(todayUSD),
		MTDUSD:            formatUSD(mtdUSD),
		PreviousPeriodUSD: formatUSD(previousUSD),
		ChangePct:         changePct(currentUSD, previousUSD),
	}, nil
}

// GetSpendTimeSeries returns daily spend points for the requested range.
func (s *Service) GetSpendTimeSeries(ctx context.Context, params QueryParams, granularity string) (TimeSeriesResult, error) {
	if granularity != "" && granularity != "daily" {
		return TimeSeriesResult{}, shared.BadRequest("invalid_granularity", "granularity must be daily")
	}

	points, err := s.clickhouse.DailyCosts(ctx, params.OrgID, params.From, params.To, params.TeamID)
	if err != nil {
		return TimeSeriesResult{}, shared.Internal("failed to query spend timeseries")
	}

	data := make([]TimeSeriesPoint, 0, len(points))
	for _, point := range points {
		data = append(data, TimeSeriesPoint{
			Date:    point.Date.UTC().Format("2006-01-02"),
			CostUSD: formatUSD(point.CostUSD),
		})
	}
	return TimeSeriesResult{Data: data}, nil
}

// GetSpendByProvider returns provider cost breakdown with percentage share.
func (s *Service) GetSpendByProvider(ctx context.Context, params QueryParams) (ProviderBreakdownResult, error) {
	items, err := s.clickhouse.SpendByProvider(ctx, params.OrgID, params.From, params.To, params.TeamID)
	if err != nil {
		return ProviderBreakdownResult{}, shared.Internal("failed to query spend by provider")
	}

	total := decimal.Zero
	for _, item := range items {
		total = total.Add(item.CostUSD)
	}

	providers := make([]ProviderBreakdownItem, 0, len(items))
	for _, item := range items {
		providers = append(providers, ProviderBreakdownItem{
			Name:    item.Provider,
			CostUSD: formatUSD(item.CostUSD),
			Pct:     sharePct(item.CostUSD, total),
		})
	}
	return ProviderBreakdownResult{Providers: providers}, nil
}

// GetSpendByTeam returns team cost breakdown enriched with PostgreSQL team names.
func (s *Service) GetSpendByTeam(ctx context.Context, params QueryParams) (TeamBreakdownResult, error) {
	items, err := s.clickhouse.SpendByTeam(ctx, params.OrgID, params.From, params.To)
	if err != nil {
		return TeamBreakdownResult{}, shared.Internal("failed to query spend by team")
	}

	teamNames, err := s.loadTeamNames(ctx, params.OrgID)
	if err != nil {
		return TeamBreakdownResult{}, err
	}

	teams := make([]TeamBreakdownItem, 0, len(items))
	for _, item := range items {
		entry := TeamBreakdownItem{
			CostUSD:    formatUSD(item.CostUSD),
			EventCount: int(item.EventCount),
		}
		if item.TeamID == nil {
			entry.ID = ""
			entry.Name = "Unassigned"
		} else {
			entry.ID = item.TeamID.String()
			entry.Name = teamNames[*item.TeamID]
			if entry.Name == "" {
				entry.Name = "Unknown team"
			}
		}
		teams = append(teams, entry)
	}
	return TeamBreakdownResult{Teams: teams}, nil
}

func (s *Service) loadTeamNames(ctx context.Context, orgID uuid.UUID) (map[uuid.UUID]string, error) {
	teams, err := s.queries.ListTeamsByOrg(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		s.log.Error("list teams failed", slog.String("org_id", orgID.String()), slog.String("error", err.Error()))
		return nil, shared.Internal("failed to load team names")
	}

	names := make(map[uuid.UUID]string, len(teams))
	for _, team := range teams {
		if !team.ID.Valid {
			continue
		}
		id, err := uuid.FromBytes(team.ID.Bytes[:])
		if err != nil {
			continue
		}
		names[id] = team.Name
	}
	return names, nil
}

func parseTimestamp(raw string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid timestamp")
}

func startOfDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func endOfDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func startOfMonth(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func formatUSD(value decimal.Decimal) string {
	return value.StringFixed(2)
}

func changePct(current, previous decimal.Decimal) float64 {
	if previous.IsZero() {
		if current.IsZero() {
			return 0
		}
		return 100
	}
	delta := current.Sub(previous)
	pct := delta.Div(previous).Mul(decimal.NewFromInt(100))
	result, _ := pct.Float64()
	return result
}

func sharePct(part, total decimal.Decimal) float64 {
	if total.IsZero() {
		return 0
	}
	pct := part.Div(total).Mul(decimal.NewFromInt(100))
	result, _ := pct.Float64()
	return result
}