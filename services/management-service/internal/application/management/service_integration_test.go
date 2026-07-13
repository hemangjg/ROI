package management_test

import (
	"context"
	"os"
	"testing"

	"github.com/ai-finops/ai-finops/packages/logger"
	mgmtapp "github.com/ai-finops/ai-finops/services/management-service/internal/application/management"
	"github.com/ai-finops/ai-finops/services/management-service/internal/infrastructure/postgres"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestManagementServiceIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, queries, err := postgres.Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer pool.Close()

	svc := mgmtapp.NewService(pool, queries, logger.New("management-service-test", "error"))

	orgID, err := svc.CreateOrg(ctx, "WS4 Org "+uuid.NewString()[:8], "")
	if err != nil {
		t.Fatalf("create org: %v", err)
	}

	adminUser, err := queries.CreateUser(ctx, sqlcgen.CreateUserParams{
		Email:        "ws4-admin-" + uuid.NewString() + "@example.com",
		PasswordHash: "test-hash",
		Name:         "WS4 Admin",
	})
	if err != nil {
		t.Fatalf("create admin user: %v", err)
	}

	actor := mgmtapp.Actor{UserID: adminUser.ID}

	_, err = queries.CreateOrgMembership(ctx, sqlcgen.CreateOrgMembershipParams{
		UserID: adminUser.ID,
		OrgID:  mustUUID(t, orgID),
		Role:   "admin",
	})
	if err != nil {
		t.Fatalf("create admin membership: %v", err)
	}

	org, err := svc.GetOrg(ctx, orgID, actor)
	if err != nil {
		t.Fatalf("get org: %v", err)
	}
	if org.ID != orgID {
		t.Fatalf("org id = %q, want %q", org.ID, orgID)
	}

	updatedName := "WS4 Updated Org"
	org, err = svc.UpdateOrg(ctx, orgID, &updatedName, nil, actor)
	if err != nil {
		t.Fatalf("update org: %v", err)
	}
	if org.Name != updatedName {
		t.Fatalf("org name = %q, want %q", org.Name, updatedName)
	}

	team, err := svc.CreateTeam(ctx, orgID, "Platform", actor)
	if err != nil {
		t.Fatalf("create team: %v", err)
	}

	teams, err := svc.ListTeams(ctx, orgID, actor)
	if err != nil || len(teams) != 1 {
		t.Fatalf("list teams: got %d teams err=%v", len(teams), err)
	}

	inviteEmail := "ws4-" + uuid.NewString() + "@example.com"
	member, err := svc.InviteUser(ctx, orgID, inviteEmail, "WS4 Member", "member", actor)
	if err != nil {
		t.Fatalf("invite user: %v", err)
	}
	if member.Email != inviteEmail {
		t.Fatalf("invited email = %q, want %q", member.Email, inviteEmail)
	}

	users, err := svc.ListUsers(ctx, orgID, actor)
	if err != nil || len(users) < 2 {
		t.Fatalf("list users: got %d users err=%v", len(users), err)
	}

	created, err := svc.CreateApiKey(ctx, orgID, "ingestion-key", actor)
	if err != nil {
		t.Fatalf("create api key: %v", err)
	}
	if created.APIKey == "" || created.KeyID == "" {
		t.Fatal("expected api key payload")
	}

	keys, err := svc.ListApiKeys(ctx, orgID, actor)
	if err != nil || len(keys) != 1 {
		t.Fatalf("list api keys: got %d keys err=%v", len(keys), err)
	}

	if err := svc.RevokeApiKey(ctx, orgID, created.KeyID, actor); err != nil {
		t.Fatalf("revoke api key: %v", err)
	}

	logs, err := svc.ListAuditLogs(ctx, orgID, "", 20, actor)
	if err != nil || len(logs.Data) == 0 {
		t.Fatalf("list audit logs: got %d entries err=%v", len(logs.Data), err)
	}

	if err := svc.DeleteTeam(ctx, orgID, team.ID, actor); err != nil {
		t.Fatalf("delete team: %v", err)
	}
}

func mustUUID(t *testing.T, raw string) pgtype.UUID {
	t.Helper()
	parsed, err := uuid.Parse(raw)
	if err != nil {
		t.Fatalf("parse uuid: %v", err)
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}
}