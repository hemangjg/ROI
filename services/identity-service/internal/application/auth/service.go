package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log/slog"
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

// TokenPair is returned by register, login, and refresh flows.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	OrgID        string `json:"org_id"`
}

// Service implements identity auth flows and validation RPCs.
type Service struct {
	pool        *pgxpool.Pool
	queries     *sqlcgen.Queries
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	accessTTL   time.Duration
	refreshTTL  time.Duration
	log         *slog.Logger
}

// NewService constructs an auth application service.
func NewService(
	pool *pgxpool.Pool,
	queries *sqlcgen.Queries,
	privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey,
	accessTTL, refreshTTL time.Duration,
	log *slog.Logger,
) *Service {
	return &Service{
		pool:       pool,
		queries:    queries,
		privateKey: privateKey,
		publicKey:  publicKey,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		log:        log,
	}
}

// Register creates a user, organization, and admin membership.
func (s *Service) Register(ctx context.Context, email, password, name, orgName string) (TokenPair, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" || name == "" || orgName == "" {
		return TokenPair{}, shared.BadRequest("invalid_request", "email, password, name, and org_name are required")
	}

	passwordHash, err := pkgauth.HashPassword(password)
	if err != nil {
		return TokenPair{}, shared.Internal("hash password")
	}

	slug := uniqueSlug(orgName)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return TokenPair{}, shared.Internal("begin transaction")
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	user, err := qtx.CreateUser(ctx, sqlcgen.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return TokenPair{}, shared.Conflict("email_taken", "email is already registered")
		}
		return TokenPair{}, shared.Internal("create user")
	}

	org, err := qtx.CreateOrganization(ctx, sqlcgen.CreateOrganizationParams{
		Name: orgName,
		Slug: slug,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return TokenPair{}, shared.Conflict("slug_taken", "organization slug is already taken")
		}
		return TokenPair{}, shared.Internal("create organization")
	}

	membership, err := qtx.CreateOrgMembership(ctx, sqlcgen.CreateOrgMembershipParams{
		UserID: user.ID,
		OrgID:  org.ID,
		Role:   "admin",
	})
	if err != nil {
		return TokenPair{}, shared.Internal("create membership")
	}

	if err := tx.Commit(ctx); err != nil {
		return TokenPair{}, shared.Internal("commit transaction")
	}

	return s.issueTokenPair(ctx, user.ID, org.ID, []string{membership.Role})
}

// Login authenticates a user and issues tokens for their primary organization.
func (s *Service) Login(ctx context.Context, email, password string) (TokenPair, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return TokenPair{}, shared.BadRequest("invalid_request", "email and password are required")
	}

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TokenPair{}, shared.Unauthorized("invalid_credentials", "invalid email or password")
		}
		return TokenPair{}, shared.Internal("lookup user")
	}

	if !pkgauth.VerifyPassword(password, user.PasswordHash) {
		return TokenPair{}, shared.Unauthorized("invalid_credentials", "invalid email or password")
	}

	memberships, err := s.queries.ListOrgMembershipsByUser(ctx, user.ID)
	if err != nil {
		return TokenPair{}, shared.Internal("list memberships")
	}
	if len(memberships) == 0 {
		return TokenPair{}, shared.Unauthorized("no_membership", "user has no organization membership")
	}

	primary := memberships[0]
	roles := []string{primary.Role}
	return s.issueTokenPair(ctx, user.ID, primary.OrgID, roles)
}

// Refresh rotates a refresh token and issues a new token pair.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return TokenPair{}, shared.BadRequest("invalid_request", "refresh_token is required")
	}

	stored, err := s.queries.GetRefreshTokenByHash(ctx, pkgauth.HashRefreshToken(refreshToken))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TokenPair{}, shared.Unauthorized("invalid_refresh_token", "refresh token is invalid")
		}
		return TokenPair{}, shared.Internal("lookup refresh token")
	}

	if stored.RevokedAt.Valid {
		return TokenPair{}, shared.Unauthorized("invalid_refresh_token", "refresh token is revoked")
	}
	if !stored.ExpiresAt.Valid || stored.ExpiresAt.Time.Before(time.Now()) {
		return TokenPair{}, shared.Unauthorized("invalid_refresh_token", "refresh token is expired")
	}

	memberships, err := s.queries.ListOrgMembershipsByUser(ctx, stored.UserID)
	if err != nil {
		return TokenPair{}, shared.Internal("list memberships")
	}
	if len(memberships) == 0 {
		return TokenPair{}, shared.Unauthorized("no_membership", "user has no organization membership")
	}

	if err := s.queries.RevokeRefreshToken(ctx, stored.ID); err != nil {
		return TokenPair{}, shared.Internal("revoke refresh token")
	}

	primary := memberships[0]
	return s.issueTokenPair(ctx, stored.UserID, primary.OrgID, []string{primary.Role})
}

// Logout revokes a refresh token.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return shared.BadRequest("invalid_request", "refresh_token is required")
	}

	stored, err := s.queries.GetRefreshTokenByHash(ctx, pkgauth.HashRefreshToken(refreshToken))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return shared.Internal("lookup refresh token")
	}

	if err := s.queries.RevokeRefreshToken(ctx, stored.ID); err != nil {
		return shared.Internal("revoke refresh token")
	}
	return nil
}

// ValidateToken validates a JWT access token for Envoy ext_authz and internal callers.
func (s *Service) ValidateToken(ctx context.Context, token string) (valid bool, userID, orgID string, roles []string, err error) {
	_ = ctx
	token = stripBearer(token)
	if token == "" {
		return false, "", "", nil, nil
	}

	claims, err := pkgauth.ParseJWT(s.publicKey, token)
	if err != nil {
		return false, "", "", nil, nil
	}

	return true, claims.UserID, claims.OrgID, claims.Roles, nil
}

// ValidateApiKey validates an API key for internal callers.
func (s *Service) ValidateApiKey(ctx context.Context, apiKey string) (valid bool, orgID, keyID string, err error) {
	apiKey = stripBearer(strings.TrimSpace(apiKey))
	if apiKey == "" {
		return false, "", "", nil
	}

	prefix, err := pkgauth.APIKeyPrefix(apiKey)
	if err != nil {
		return false, "", "", nil
	}

	row, err := s.queries.GetApiKeyByPrefix(ctx, prefix)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, "", "", nil
		}
		return false, "", "", shared.Internal("lookup api key")
	}

	ok, err := pkgauth.VerifyAPIKey(apiKey, row.KeyHash)
	if err != nil || !ok {
		return false, "", "", nil
	}

	return true, pgUUIDString(row.OrgID), pgUUIDString(row.ID), nil
}

func (s *Service) issueTokenPair(ctx context.Context, userID, orgID pgtype.UUID, roles []string) (TokenPair, error) {
	accessToken, err := pkgauth.SignJWT(s.privateKey, pkgauth.NewAccessClaims(
		pgUUIDString(userID),
		pgUUIDString(orgID),
		roles,
		s.accessTTL,
	))
	if err != nil {
		return TokenPair{}, shared.Internal("sign access token")
	}

	refreshPlain, refreshHash, err := pkgauth.GenerateRefreshToken()
	if err != nil {
		return TokenPair{}, shared.Internal("generate refresh token")
	}

	_, err = s.queries.CreateRefreshToken(ctx, sqlcgen.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: refreshHash,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(s.refreshTTL), Valid: true},
	})
	if err != nil {
		return TokenPair{}, shared.Internal("store refresh token")
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshPlain,
		OrgID:        pgUUIDString(orgID),
	}, nil
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

func stripBearer(token string) string {
	token = strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return strings.TrimSpace(token[7:])
	}
	return token
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

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}