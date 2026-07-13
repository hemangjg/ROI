package analytics_test

import (
	"testing"
	"time"

	analyticsapp "github.com/ai-finops/ai-finops/services/analytics-service/internal/application/analytics"
	"github.com/google/uuid"
)

func TestParseQueryParamsDefaults(t *testing.T) {
	orgID := uuid.New().String()

	params, err := analyticsapp.ParseQueryParams(orgID, "", "", "")
	if err != nil {
		t.Fatalf("ParseQueryParams() error = %v", err)
	}

	if params.OrgID.String() != orgID {
		t.Fatalf("org_id = %s, want %s", params.OrgID, orgID)
	}
	if params.TeamID != nil {
		t.Fatal("expected nil team_id")
	}

	now := time.Now().UTC()
	expectedTo := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	expectedFrom := expectedTo.AddDate(0, 0, -29)

	if !params.To.Equal(expectedTo) {
		t.Fatalf("to = %s, want %s", params.To, expectedTo)
	}
	if !params.From.Equal(expectedFrom) {
		t.Fatalf("from = %s, want %s", params.From, expectedFrom)
	}
}

func TestParseQueryParamsInvalidOrgID(t *testing.T) {
	_, err := analyticsapp.ParseQueryParams("not-a-uuid", "", "", "")
	if err == nil {
		t.Fatal("expected error for invalid org_id")
	}
}

func TestParseQueryParamsInvalidRange(t *testing.T) {
	orgID := uuid.New().String()
	_, err := analyticsapp.ParseQueryParams(orgID, "2026-07-10", "2026-07-01", "")
	if err == nil {
		t.Fatal("expected error for invalid date range")
	}
}

func TestParseQueryParamsTeamFilter(t *testing.T) {
	orgID := uuid.New().String()
	teamID := uuid.New().String()

	params, err := analyticsapp.ParseQueryParams(orgID, "", "", teamID)
	if err != nil {
		t.Fatalf("ParseQueryParams() error = %v", err)
	}
	if params.TeamID == nil || params.TeamID.String() != teamID {
		t.Fatalf("team_id = %v, want %s", params.TeamID, teamID)
	}
}