package grpcapi

import (
	"context"

	mgmtapp "github.com/ai-finops/ai-finops/services/management-service/internal/application/management"
	managementv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/management/v1"
)

// ManagementServer implements management.v1.ManagementService.
type ManagementServer struct {
	managementv1.UnimplementedManagementServiceServer
	svc *mgmtapp.Service
}

// NewManagementServer constructs the gRPC management service.
func NewManagementServer(svc *mgmtapp.Service) *ManagementServer {
	return &ManagementServer{svc: svc}
}

// CreateOrg creates a new organization.
func (s *ManagementServer) CreateOrg(ctx context.Context, req *managementv1.CreateOrgRequest) (*managementv1.CreateOrgResponse, error) {
	orgID, err := s.svc.CreateOrg(ctx, req.GetName(), req.GetSlug())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &managementv1.CreateOrgResponse{OrgId: orgID}, nil
}

// CreateApiKey creates a new API key for an organization.
func (s *ManagementServer) CreateApiKey(ctx context.Context, req *managementv1.CreateApiKeyRequest) (*managementv1.CreateApiKeyResponse, error) {
	created, err := s.svc.CreateApiKey(ctx, req.GetOrgId(), req.GetName(), mgmtapp.Actor{})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &managementv1.CreateApiKeyResponse{
		KeyId:     created.KeyID,
		ApiKey:    created.APIKey,
		KeyPrefix: created.KeyPrefix,
	}, nil
}