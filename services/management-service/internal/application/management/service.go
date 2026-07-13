package management

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"regexp"
	"strings"
	"time"

	pkgauth "github.com/ai-finops/ai-finops/packages/auth"
	"github.com/ai-finops/ai-finops/packages/shared"
	sqlcgen "github.com/ai-finops/ai-finops/db/atlas/sql/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var slugSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

const (
	roleAdmin  = "admin"
	roleMember = "member"
	roleViewer = "viewer"
)

// Actor identifies the caller for optional RBAC and audit logging.
type Actor struct {
	UserID pgtype.UUID
	IP     *netip.Addr
}

// Organization is the REST representation of an org.
type Organization struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Timezone  string `json:"timezone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Team is the REST representation of a team.
type Team struct {
	ID        string `json:"id"`
	OrgID     string `json:"org_id,omitempty"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// OrgMember is a user within an organization.
type OrgMember struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}

// ApiKeyCreated is returned once when a key is created.
type ApiKeyCreated struct {
	KeyID     string `json:"key_id"`
	APIKey    string `json:"api_key"`
	KeyPrefix string `json:"key_prefix"`
}

// ApiKey is a masked API key listing entry.
type ApiKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	KeyPrefix string `json:"key_prefix"`
	CreatedAt string `json:"created_at"`
}

// AuditLogEntry is a single audit log row.
type AuditLogEntry struct {
	ID           string          `json:"id"`
	ActorID      string          `json:"actor_id,omitempty"`
	ActorType    string          `json:"actor_type"`
	Action       string          `json:"action"`
	ResourceType string          `json:"resource_type"`
	ResourceID   string          `json:"resource_id,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	CreatedAt    string          `json:"created_at"`
}

// AuditLogPage is a cursor-paginated audit log response.
type AuditLogPage struct {
	Data       []AuditLogEntry `json:"data"`
	NextCursor string          `json:"next_cursor,omitempty"`
}

// Service implements management domain operations.
type Service struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
	log     *slog.Logger
}

// NewService constructs a management application service.
func NewService(pool *pgxpool.Pool, queries *sqlcgen.Queries, log *slog.Logger) *Service {
	return &Service{pool: pool, queries: queries, log: log}
}

// GetOrg returns organization details.
func (s *Service) GetOrg(ctx context.Context, orgID string, actor Actor) (Organization, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return Organization{}, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleViewer); err != nil {
		return Organization{}, err
	}

	row, err := s.queries.GetOrganizationByID(ctx, orgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Organization{}, shared.NotFound("org_not_found", "organization not found")
		}
		return Organization{}, shared.Internal("get organization")
	}

	return toOrganization(row), nil
}

// UpdateOrg patches organization settings.
func (s *Service) UpdateOrg(ctx context.Context, orgID string, name, timezone *string, actor Actor) (Organization, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return Organization{}, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return Organization{}, err
	}
	if name == nil && timezone == nil {
		return Organization{}, shared.BadRequest("invalid_request", "at least one of name or timezone is required")
	}

	row, err := s.queries.UpdateOrganization(ctx, sqlcgen.UpdateOrganizationParams{
		ID:       orgUUID,
		Name:     optionalText(name),
		Timezone: optionalText(timezone),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Organization{}, shared.NotFound("org_not_found", "organization not found")
		}
		return Organization{}, shared.Internal("update organization")
	}

	s.audit(ctx, orgUUID, actor, "org.updated", "organization", orgID, nil)
	return toOrganization(row), nil
}

// CreateOrg creates a new organization (gRPC and internal use).
func (s *Service) CreateOrg(ctx context.Context, name, slug string) (string, error) {
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)
	if name == "" {
		return "", shared.BadRequest("invalid_request", "name is required")
	}
	if slug == "" {
		slug = uniqueSlug(name)
	}

	row, err := s.queries.CreateOrganization(ctx, sqlcgen.CreateOrganizationParams{
		Name:    name,
		Slug:    slug,
		Column3: "UTC",
	})
	if err != nil {
		if isUniqueViolation(err) {
			return "", shared.Conflict("slug_taken", "organization slug is already taken")
		}
		return "", shared.Internal("create organization")
	}

	return pgUUIDString(row.ID), nil
}

// CreateTeam creates a team within an organization.
func (s *Service) CreateTeam(ctx context.Context, orgID, name string, actor Actor) (Team, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return Team{}, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return Team{}, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return Team{}, shared.BadRequest("invalid_request", "name is required")
	}

	row, err := s.queries.CreateTeam(ctx, sqlcgen.CreateTeamParams{
		OrgID: orgUUID,
		Name:  name,
	})
	if err != nil {
		return Team{}, shared.Internal("create team")
	}

	teamID := pgUUIDString(row.ID)
	s.audit(ctx, orgUUID, actor, "team.created", "team", teamID, map[string]string{"name": name})
	return toTeam(row), nil
}

// ListTeams returns active teams for an organization.
func (s *Service) ListTeams(ctx context.Context, orgID string, actor Actor) ([]Team, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return nil, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleViewer); err != nil {
		return nil, err
	}

	rows, err := s.queries.ListTeamsByOrg(ctx, orgUUID)
	if err != nil {
		return nil, shared.Internal("list teams")
	}

	teams := make([]Team, 0, len(rows))
	for _, row := range rows {
		teams = append(teams, toTeam(row))
	}
	return teams, nil
}

// DeleteTeam soft-deletes a team.
func (s *Service) DeleteTeam(ctx context.Context, orgID, teamID string, actor Actor) error {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return err
	}
	teamUUID, err := parseUUID(teamID)
	if err != nil {
		return err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return err
	}

	if _, err := s.queries.GetTeamByID(ctx, sqlcgen.GetTeamByIDParams{
		ID:    teamUUID,
		OrgID: orgUUID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.NotFound("team_not_found", "team not found")
		}
		return shared.Internal("lookup team")
	}

	if err := s.queries.SoftDeleteTeam(ctx, sqlcgen.SoftDeleteTeamParams{
		ID:    teamUUID,
		OrgID: orgUUID,
	}); err != nil {
		return shared.Internal("delete team")
	}

	s.audit(ctx, orgUUID, actor, "team.deleted", "team", teamID, nil)
	return nil
}

// InviteUser adds a user to the organization, creating the user if needed.
func (s *Service) InviteUser(ctx context.Context, orgID, email, name, role string, actor Actor) (OrgMember, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return OrgMember{}, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return OrgMember{}, err
	}

	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)
	role = strings.TrimSpace(role)
	if email == "" || name == "" {
		return OrgMember{}, shared.BadRequest("invalid_request", "email and name are required")
	}
	if role == "" {
		role = roleMember
	}
	if !isValidRole(role) {
		return OrgMember{}, shared.BadRequest("invalid_role", "role must be admin, member, or viewer")
	}

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return OrgMember{}, shared.Internal("lookup user")
	}

	if errors.Is(err, pgx.ErrNoRows) {
		passwordHash, hashErr := pkgauth.HashPassword(randomInvitePassword())
		if hashErr != nil {
			return OrgMember{}, shared.Internal("hash invite password")
		}
		user, err = s.queries.CreateUser(ctx, sqlcgen.CreateUserParams{
			Email:        email,
			PasswordHash: passwordHash,
			Name:         name,
		})
		if err != nil {
			if isUniqueViolation(err) {
				user, err = s.queries.GetUserByEmail(ctx, email)
				if err != nil {
					return OrgMember{}, shared.Internal("lookup user after conflict")
				}
			} else {
				return OrgMember{}, shared.Internal("create user")
			}
		}
	}

	membership, err := s.queries.CreateOrgMembership(ctx, sqlcgen.CreateOrgMembershipParams{
		UserID: user.ID,
		OrgID:  orgUUID,
		Role:   role,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return OrgMember{}, shared.Conflict("already_member", "user is already a member of this organization")
		}
		return OrgMember{}, shared.Internal("create membership")
	}

	userID := pgUUIDString(user.ID)
	s.audit(ctx, orgUUID, actor, "user.invited", "user", userID, map[string]string{
		"email": email,
		"role":  role,
	})

	return OrgMember{
		ID:       userID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     membership.Role,
		JoinedAt: formatTime(membership.CreatedAt),
	}, nil
}

// ListUsers returns organization members.
func (s *Service) ListUsers(ctx context.Context, orgID string, actor Actor) ([]OrgMember, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return nil, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return nil, err
	}

	rows, err := s.queries.ListOrgMembersByOrg(ctx, orgUUID)
	if err != nil {
		return nil, shared.Internal("list members")
	}

	members := make([]OrgMember, 0, len(rows))
	for _, row := range rows {
		members = append(members, OrgMember{
			ID:       pgUUIDString(row.ID),
			Email:    row.Email,
			Name:     row.Name,
			Role:     row.Role,
			JoinedAt: formatTime(row.JoinedAt),
		})
	}
	return members, nil
}

// CreateApiKey creates a new API key for an organization.
func (s *Service) CreateApiKey(ctx context.Context, orgID, name string, actor Actor) (ApiKeyCreated, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return ApiKeyCreated{}, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return ApiKeyCreated{}, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return ApiKeyCreated{}, shared.BadRequest("invalid_request", "name is required")
	}

	plaintext, hash, err := pkgauth.GenerateAPIKey()
	if err != nil {
		return ApiKeyCreated{}, shared.Internal("generate api key")
	}

	prefix, err := pkgauth.APIKeyPrefix(plaintext)
	if err != nil {
		return ApiKeyCreated{}, shared.Internal("derive api key prefix")
	}

	row, err := s.queries.CreateApiKey(ctx, sqlcgen.CreateApiKeyParams{
		OrgID:     orgUUID,
		Name:      name,
		KeyPrefix: prefix,
		KeyHash:   hash,
		Scopes:    []string{"ingest"},
		CreatedBy: actor.UserID,
	})
	if err != nil {
		return ApiKeyCreated{}, shared.Internal("create api key")
	}

	keyID := pgUUIDString(row.ID)
	s.audit(ctx, orgUUID, actor, "api_key.created", "api_key", keyID, map[string]string{"name": name})

	return ApiKeyCreated{
		KeyID:     keyID,
		APIKey:    plaintext,
		KeyPrefix: prefix,
	}, nil
}

// ListApiKeys returns active API keys for an organization.
func (s *Service) ListApiKeys(ctx context.Context, orgID string, actor Actor) ([]ApiKey, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return nil, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return nil, err
	}

	rows, err := s.queries.ListApiKeysByOrg(ctx, orgUUID)
	if err != nil {
		return nil, shared.Internal("list api keys")
	}

	keys := make([]ApiKey, 0, len(rows))
	for _, row := range rows {
		keys = append(keys, ApiKey{
			ID:        pgUUIDString(row.ID),
			Name:      row.Name,
			KeyPrefix: row.KeyPrefix,
			CreatedAt: formatTime(row.CreatedAt),
		})
	}
	return keys, nil
}

// RevokeApiKey revokes an API key.
func (s *Service) RevokeApiKey(ctx context.Context, orgID, keyID string, actor Actor) error {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return err
	}
	keyUUID, err := parseUUID(keyID)
	if err != nil {
		return err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return err
	}

	if _, err := s.queries.GetApiKeyByID(ctx, sqlcgen.GetApiKeyByIDParams{
		ID:    keyUUID,
		OrgID: orgUUID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.NotFound("api_key_not_found", "api key not found")
		}
		return shared.Internal("lookup api key")
	}

	if err := s.queries.RevokeApiKey(ctx, sqlcgen.RevokeApiKeyParams{
		ID:    keyUUID,
		OrgID: orgUUID,
	}); err != nil {
		return shared.Internal("revoke api key")
	}

	s.audit(ctx, orgUUID, actor, "api_key.revoked", "api_key", keyID, nil)
	return nil
}

// ListAuditLogs returns paginated audit logs for an organization.
func (s *Service) ListAuditLogs(ctx context.Context, orgID, cursor string, limit int, actor Actor) (AuditLogPage, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return AuditLogPage{}, err
	}
	if err := s.requireRole(ctx, actor, orgUUID, roleAdmin); err != nil {
		return AuditLogPage{}, err
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	var cursorTS pgtype.Timestamptz
	if strings.TrimSpace(cursor) != "" {
		parsed, parseErr := time.Parse(time.RFC3339Nano, cursor)
		if parseErr != nil {
			return AuditLogPage{}, shared.BadRequest("invalid_cursor", "cursor must be RFC3339Nano timestamp")
		}
		cursorTS = pgtype.Timestamptz{Time: parsed, Valid: true}
	}

	rows, err := s.queries.ListAuditLogsByOrg(ctx, sqlcgen.ListAuditLogsByOrgParams{
		OrgID:           orgUUID,
		CursorCreatedAt: cursorTS,
		PageLimit:       int32(limit),
	})
	if err != nil {
		return AuditLogPage{}, shared.Internal("list audit logs")
	}

	page := AuditLogPage{Data: make([]AuditLogEntry, 0, len(rows))}
	for _, row := range rows {
		entry := AuditLogEntry{
			ID:           pgUUIDString(row.ID),
			ActorType:    row.ActorType,
			Action:       row.Action,
			ResourceType: row.ResourceType,
			CreatedAt:    formatTime(row.CreatedAt),
		}
		if row.ActorID.Valid {
			entry.ActorID = pgUUIDString(row.ActorID)
		}
		if row.ResourceID.Valid {
			entry.ResourceID = row.ResourceID.String
		}
		if len(row.Metadata) > 0 {
			entry.Metadata = json.RawMessage(row.Metadata)
		}
		page.Data = append(page.Data, entry)
	}

	if len(rows) == limit {
		last := rows[len(rows)-1]
		page.NextCursor = last.CreatedAt.Time.Format(time.RFC3339Nano)
	}

	return page, nil
}

func (s *Service) requireRole(ctx context.Context, actor Actor, orgID pgtype.UUID, minRole string) error {
	if !actor.UserID.Valid {
		return nil
	}

	membership, err := s.queries.GetOrgMembership(ctx, sqlcgen.GetOrgMembershipParams{
		UserID: actor.UserID,
		OrgID:  orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.Forbidden("not_member", "user is not a member of this organization")
		}
		return shared.Internal("lookup membership")
	}

	if roleRank(membership.Role) < roleRank(minRole) {
		return shared.Forbidden("insufficient_role", "insufficient permissions for this operation")
	}
	return nil
}

func (s *Service) audit(ctx context.Context, orgID pgtype.UUID, actor Actor, action, resourceType, resourceID string, metadata any) {
	var metaBytes []byte
	if metadata != nil {
		metaBytes, _ = json.Marshal(metadata)
	}

	actorType := "system"
	if actor.UserID.Valid {
		actorType = "user"
	}

	_, err := s.queries.InsertAuditLog(ctx, sqlcgen.InsertAuditLogParams{
		OrgID:        orgID,
		ActorID:      actor.UserID,
		ActorType:    actorType,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   pgtype.Text{String: resourceID, Valid: resourceID != ""},
		Metadata:     metaBytes,
		IpAddress:    actor.IP,
	})
	if err != nil {
		s.log.Warn("failed to write audit log",
			slog.String("action", action),
			slog.String("resource_type", resourceType),
			slog.String("error", err.Error()),
		)
	}
}

func roleRank(role string) int {
	switch role {
	case roleAdmin:
		return 3
	case roleMember:
		return 2
	case roleViewer:
		return 1
	default:
		return 0
	}
}

func isValidRole(role string) bool {
	return role == roleAdmin || role == roleMember || role == roleViewer
}

func parseUUID(raw string) (pgtype.UUID, error) {
	raw = strings.TrimSpace(raw)
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return pgtype.UUID{}, shared.BadRequest("invalid_uuid", "invalid id format")
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

func pgUUIDString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	u, err := uuid.FromBytes(id.Bytes[:])
	if err != nil {
		return ""
	}
	return u.String()
}

func formatTime(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.UTC().Format(time.RFC3339)
}

func toOrganization(row sqlcgen.Organization) Organization {
	return Organization{
		ID:        pgUUIDString(row.ID),
		Name:      row.Name,
		Slug:      row.Slug,
		Timezone:  row.Timezone,
		CreatedAt: formatTime(row.CreatedAt),
		UpdatedAt: formatTime(row.UpdatedAt),
	}
}

func toTeam(row sqlcgen.Team) Team {
	return Team{
		ID:        pgUUIDString(row.ID),
		OrgID:     pgUUIDString(row.OrgID),
		Name:      row.Name,
		CreatedAt: formatTime(row.CreatedAt),
	}
}

func optionalText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: trimmed, Valid: true}
}

func uniqueSlug(orgName string) string {
	base := slugify(orgName)
	return fmt.Sprintf("%s-%s", base, uuid.NewString()[:8])
}

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = slugSanitizer.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "org"
	}
	return s
}

func randomInvitePassword() string {
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		return uuid.NewString()
	}
	return base64.RawURLEncoding.EncodeToString(raw)
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}