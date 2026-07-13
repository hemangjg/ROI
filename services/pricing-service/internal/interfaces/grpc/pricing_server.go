package grpcapi

import (
	"context"
	"time"

	pricingapp "github.com/ai-finops/ai-finops/services/pricing-service/internal/application/pricing"
	pricingv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/pricing/v1"
)

// PricingServer implements pricing.v1.PricingService.
type PricingServer struct {
	pricingv1.UnimplementedPricingServiceServer
	svc *pricingapp.Service
}

// NewPricingServer constructs the gRPC pricing service.
func NewPricingServer(svc *pricingapp.Service) *PricingServer {
	return &PricingServer{svc: svc}
}

// CalculateCost computes event cost from the versioned pricing catalog.
func (s *PricingServer) CalculateCost(ctx context.Context, req *pricingv1.CalculateCostRequest) (*pricingv1.CalculateCostResponse, error) {
	occurredAt := time.Now().UTC()
	if req.GetOccurredAt() != nil {
		occurredAt = req.GetOccurredAt().AsTime().UTC()
	}

	result, err := s.svc.CalculateCost(
		ctx,
		req.GetProvider(),
		req.GetModel(),
		req.GetInputTokens(),
		req.GetOutputTokens(),
		req.GetCacheReadTokens(),
		req.GetCacheWriteTokens(),
		occurredAt,
	)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &pricingv1.CalculateCostResponse{
		CostUsd:   result.CostUSD,
		IsPriced:  result.IsPriced,
		ModelId:   result.ModelID,
	}, nil
}

// SyncPricingCatalog refreshes provider pricing from embedded plugins.
func (s *PricingServer) SyncPricingCatalog(ctx context.Context, _ *pricingv1.SyncPricingCatalogRequest) (*pricingv1.SyncPricingCatalogResponse, error) {
	updated, err := s.svc.SyncPricingCatalog(ctx)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pricingv1.SyncPricingCatalogResponse{ModelsUpdated: updated}, nil
}