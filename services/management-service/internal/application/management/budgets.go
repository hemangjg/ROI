package management

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Budget is the REST representation of a team monthly budget.
type Budget struct {
	ID                string `json:"id"`
	OrgID             string `json:"org_id"`
	TeamID            string `json:"team_id"`
	TeamName          string `json:"team_name,omitempty"`
	AmountUSD         string `json:"amount_usd"`
	SoftThresholdPct  int32  `json:"soft_threshold_pct"`
	PeriodStart       string `json:"period_start"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// BudgetAlert is a soft-limit notification for a budget.
type BudgetAlert struct {
	ID           string `json:"id"`
	BudgetID     string `json:"budget_id"`
	OrgID        string `json:"org_id"`
	TeamID       string `json:"team_id"`
	ThresholdPct int32  `json:"threshold_pct"`
	SpendUSD     string `json:"spend_usd"`
	BudgetUSD    string `json:"budget_usd"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	CreatedAt    string `json:"created_at"`
}

// BudgetStatus combines budget config with reported spend.
type BudgetStatus struct {
	Budget
	SpendUSD     string  `json:"spend_usd"`
	UsagePct     float64 `json:"usage_pct"`
	OverSoft     bool    `json:"over_soft_threshold"`
	AlertFired   bool    `json:"alert_fired"`
	Alert        *BudgetAlert `json:"alert,omitempty"`
}

// SoftThresholdBreached returns true when spend reaches the soft limit.
func SoftThresholdBreached(spendUSD, budgetUSD float64, thresholdPct int32) bool {
	if budgetUSD <= 0 || thresholdPct <= 0 {
		return false
	}
	return spendUSD*100 >= budgetUSD*float64(thresholdPct)
}

// UsagePercent returns spend/budget * 100.
func UsagePercent(spendUSD, budgetUSD float64) float64 {
	if budgetUSD <= 0 {
		return 0
	}
	return (spendUSD / budgetUSD) * 100
}

// UpsertBudget creates or updates a team's monthly budget (admin).
func (s *Service) UpsertBudget(ctx context.Context, orgID, teamID, amountUSD string, softThresholdPct int32, periodStart string, actor Actor) (Budget, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return Budget{}, err
	}
	teamUUID, err := parseUUID(teamID)
	if err != nil {
		return Budget{}, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return Budget{}, err
	}

	if softThresholdPct <= 0 {
		softThresholdPct = 80
	}
	if softThresholdPct > 100 {
		return Budget{}, shared.BadRequest("invalid_threshold", "soft_threshold_pct must be between 1 and 100")
	}

	amount, err := parseUSD(amountUSD)
	if err != nil {
		return Budget{}, err
	}
	period, err := parsePeriodStart(periodStart)
	if err != nil {
		return Budget{}, err
	}

	if _, err := s.queries.GetTeamByID(ctx, sqlcgen.GetTeamByIDParams{ID: teamUUID, OrgID: orgUUID}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Budget{}, shared.NotFound("team_not_found", "team not found")
		}
		return Budget{}, shared.Internal("lookup team")
	}

	row, err := s.queries.UpsertBudget(ctx, sqlcgen.UpsertBudgetParams{
		OrgID:             orgUUID,
		TeamID:            teamUUID,
		AmountUsd:         amount,
		SoftThresholdPct:  softThresholdPct,
		PeriodStart:       period,
		CreatedBy:         actor.UserID,
	})
	if err != nil {
		return Budget{}, shared.Internal("upsert budget")
	}

	s.audit(ctx, orgUUID, actor, "budget.upsert", "budget", pgUUIDString(row.ID), map[string]any{
		"team_id": teamID,
		"amount_usd": amountUSD,
		"period_start": formatDate(period),
	})

	return toBudget(row, ""), nil
}

// ListBudgets lists budgets for an org period (default: current month UTC).
func (s *Service) ListBudgets(ctx context.Context, orgID, periodStart string, actor Actor) ([]Budget, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return nil, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleViewer); err != nil {
		return nil, err
	}
	period, err := parsePeriodStart(periodStart)
	if err != nil {
		return nil, err
	}
	rows, err := s.queries.ListBudgetsByOrgPeriod(ctx, sqlcgen.ListBudgetsByOrgPeriodParams{
		OrgID:       orgUUID,
		PeriodStart: period,
	})
	if err != nil {
		return nil, shared.Internal("list budgets")
	}
	out := make([]Budget, 0, len(rows))
	for _, row := range rows {
		out = append(out, Budget{
			ID:               pgUUIDString(row.ID),
			OrgID:            pgUUIDString(row.OrgID),
			TeamID:           pgUUIDString(row.TeamID),
			TeamName:         row.TeamName,
			AmountUSD:        numericString(row.AmountUsd),
			SoftThresholdPct: row.SoftThresholdPct,
			PeriodStart:      formatDate(row.PeriodStart),
			CreatedAt:        formatTime(row.CreatedAt),
			UpdatedAt:        formatTime(row.UpdatedAt),
		})
	}
	return out, nil
}

// EvaluateBudget soft-checks spend against a team budget and records an alert if breached.
func (s *Service) EvaluateBudget(ctx context.Context, orgID, teamID, periodStart, spendUSD string, actor Actor) (BudgetStatus, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return BudgetStatus{}, err
	}
	teamUUID, err := parseUUID(teamID)
	if err != nil {
		return BudgetStatus{}, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleMember); err != nil {
		return BudgetStatus{}, err
	}
	period, err := parsePeriodStart(periodStart)
	if err != nil {
		return BudgetStatus{}, err
	}
	spend, err := parseUSD(spendUSD)
	if err != nil {
		return BudgetStatus{}, err
	}

	budgetRow, err := s.queries.GetBudgetByTeamPeriod(ctx, sqlcgen.GetBudgetByTeamPeriodParams{
		OrgID:       orgUUID,
		TeamID:      teamUUID,
		PeriodStart: period,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BudgetStatus{}, shared.NotFound("budget_not_found", "no budget for team/period")
		}
		return BudgetStatus{}, shared.Internal("get budget")
	}

	spendF, _ := spend.Float64Value()
	budgetF, _ := budgetRow.AmountUsd.Float64Value()
	spendVal := spendF.Float64
	budgetVal := budgetF.Float64
	usage := UsagePercent(spendVal, budgetVal)
	over := SoftThresholdBreached(spendVal, budgetVal, budgetRow.SoftThresholdPct)

	status := BudgetStatus{
		Budget:   toBudget(budgetRow, ""),
		SpendUSD: numericString(spend),
		UsagePct: usage,
		OverSoft: over,
	}

	if over {
		msg := fmt.Sprintf(
			"Team budget soft limit reached: spent $%s of $%s (%d%% threshold)",
			numericString(spend),
			numericString(budgetRow.AmountUsd),
			budgetRow.SoftThresholdPct,
		)
		alertRow, alertErr := s.queries.InsertBudgetAlert(ctx, sqlcgen.InsertBudgetAlertParams{
			BudgetID:     budgetRow.ID,
			OrgID:        orgUUID,
			TeamID:       teamUUID,
			ThresholdPct: budgetRow.SoftThresholdPct,
			SpendUsd:     spend,
			BudgetUsd:    budgetRow.AmountUsd,
			Message:      msg,
		})
		if alertErr == nil {
			a := toBudgetAlert(alertRow)
			status.AlertFired = true
			status.Alert = &a
			s.audit(ctx, orgUUID, actor, "budget.alert", "budget_alert", a.ID, map[string]any{
				"team_id": teamID,
				"threshold_pct": budgetRow.SoftThresholdPct,
			})
		} else if !errors.Is(alertErr, pgx.ErrNoRows) {
			// ON CONFLICT DO NOTHING returns no rows via QueryRow — treat as already alerted.
			if !isNoRows(alertErr) {
				s.log.Warn("insert budget alert failed", slog.String("error", alertErr.Error()))
			}
		}
	}

	return status, nil
}

// ListBudgetAlerts returns recent budget alerts for an org.
func (s *Service) ListBudgetAlerts(ctx context.Context, orgID, status string, limit int, actor Actor) ([]BudgetAlert, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return nil, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleViewer); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var statusFilter pgtype.Text
	status = strings.TrimSpace(status)
	if status != "" {
		if status != "open" && status != "acknowledged" {
			return nil, shared.BadRequest("invalid_status", "status must be open or acknowledged")
		}
		statusFilter = pgtype.Text{String: status, Valid: true}
	}
	rows, err := s.queries.ListBudgetAlertsByOrg(ctx, sqlcgen.ListBudgetAlertsByOrgParams{
		OrgID:        orgUUID,
		StatusFilter: statusFilter,
		PageLimit:    int32(limit),
	})
	if err != nil {
		return nil, shared.Internal("list budget alerts")
	}
	out := make([]BudgetAlert, 0, len(rows))
	for _, row := range rows {
		out = append(out, toBudgetAlert(row))
	}
	return out, nil
}

// AcknowledgeBudgetAlert marks an open alert as acknowledged (admin/member).
func (s *Service) AcknowledgeBudgetAlert(ctx context.Context, orgID, alertID string, actor Actor) (BudgetAlert, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return BudgetAlert{}, err
	}
	alertUUID, err := parseUUID(alertID)
	if err != nil {
		return BudgetAlert{}, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleMember); err != nil {
		return BudgetAlert{}, err
	}
	row, err := s.queries.AcknowledgeBudgetAlert(ctx, sqlcgen.AcknowledgeBudgetAlertParams{
		ID:    alertUUID,
		OrgID: orgUUID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BudgetAlert{}, shared.NotFound("alert_not_found", "open alert not found")
		}
		return BudgetAlert{}, shared.Internal("acknowledge alert")
	}
	s.audit(ctx, orgUUID, actor, "budget.alert.ack", "budget_alert", alertID, nil)
	return toBudgetAlert(row), nil
}

func parseUSD(raw string) (pgtype.Numeric, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return pgtype.Numeric{}, shared.BadRequest("invalid_amount", "amount_usd is required")
	}
	var n pgtype.Numeric
	if err := n.Scan(raw); err != nil {
		return pgtype.Numeric{}, shared.BadRequest("invalid_amount", "amount_usd must be a number")
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid || f.Float64 <= 0 {
		return pgtype.Numeric{}, shared.BadRequest("invalid_amount", "amount_usd must be greater than 0")
	}
	return n, nil
}

func parsePeriodStart(raw string) (pgtype.Date, error) {
	raw = strings.TrimSpace(raw)
	var t time.Time
	var err error
	if raw == "" {
		now := time.Now().UTC()
		t = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else if len(raw) == 7 {
		// YYYY-MM
		t, err = time.Parse("2006-01", raw)
		if err != nil {
			return pgtype.Date{}, shared.BadRequest("invalid_period", "period must be YYYY-MM or YYYY-MM-DD")
		}
		t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		t, err = time.Parse("2006-01-02", raw)
		if err != nil {
			return pgtype.Date{}, shared.BadRequest("invalid_period", "period must be YYYY-MM or YYYY-MM-DD")
		}
		t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

func formatDate(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

func numericString(n pgtype.Numeric) string {
	if !n.Valid {
		return "0"
	}
	// Prefer fixed 2-decimal money display when possible.
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return "0"
	}
	return fmt.Sprintf("%.2f", f.Float64)
}

func toBudget(row sqlcgen.Budget, teamName string) Budget {
	return Budget{
		ID:               pgUUIDString(row.ID),
		OrgID:            pgUUIDString(row.OrgID),
		TeamID:           pgUUIDString(row.TeamID),
		TeamName:         teamName,
		AmountUSD:        numericString(row.AmountUsd),
		SoftThresholdPct: row.SoftThresholdPct,
		PeriodStart:      formatDate(row.PeriodStart),
		CreatedAt:        formatTime(row.CreatedAt),
		UpdatedAt:        formatTime(row.UpdatedAt),
	}
}

func toBudgetAlert(row sqlcgen.BudgetAlert) BudgetAlert {
	return BudgetAlert{
		ID:           pgUUIDString(row.ID),
		BudgetID:     pgUUIDString(row.BudgetID),
		OrgID:        pgUUIDString(row.OrgID),
		TeamID:       pgUUIDString(row.TeamID),
		ThresholdPct: row.ThresholdPct,
		SpendUSD:     numericString(row.SpendUsd),
		BudgetUSD:    numericString(row.BudgetUsd),
		Status:       row.Status,
		Message:      row.Message,
		CreatedAt:    formatTime(row.CreatedAt),
	}
}

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows) || strings.Contains(strings.ToLower(err.Error()), "no rows")
}
