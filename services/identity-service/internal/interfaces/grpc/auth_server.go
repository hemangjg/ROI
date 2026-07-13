package grpcapi

import (
	"context"

	authapp "github.com/ai-finops/ai-finops/services/identity-service/internal/application/auth"
	authv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/auth/v1"
)

// AuthServer implements auth.v1.AuthService.
type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	svc *authapp.Service
}

// NewAuthServer constructs the gRPC auth service.
func NewAuthServer(svc *authapp.Service) *AuthServer {
	return &AuthServer{svc: svc}
}

// ValidateToken validates a JWT access token.
func (s *AuthServer) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	valid, userID, orgID, roles, err := s.svc.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &authv1.ValidateTokenResponse{
		Valid:  valid,
		UserId: userID,
		OrgId:  orgID,
		Roles:  roles,
	}, nil
}

// ValidateApiKey validates an API key.
func (s *AuthServer) ValidateApiKey(ctx context.Context, req *authv1.ValidateApiKeyRequest) (*authv1.ValidateApiKeyResponse, error) {
	valid, orgID, keyID, err := s.svc.ValidateApiKey(ctx, req.GetApiKey())
	if err != nil {
		return nil, err
	}

	return &authv1.ValidateApiKeyResponse{
		Valid: valid,
		OrgId: orgID,
		KeyId: keyID,
	}, nil
}