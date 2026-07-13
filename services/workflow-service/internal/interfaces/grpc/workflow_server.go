package grpcapi

import (
	"context"
	"time"

	workflowapp "github.com/ai-finops/ai-finops/services/workflow-service/internal/application/workflow"
	workflowv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/workflow/v1"
)

// WorkflowServer implements workflow.v1.WorkflowService.
type WorkflowServer struct {
	workflowv1.UnimplementedWorkflowServiceServer
	svc *workflowapp.Service
}

// NewWorkflowServer constructs the gRPC workflow service.
func NewWorkflowServer(svc *workflowapp.Service) *WorkflowServer {
	return &WorkflowServer{svc: svc}
}

// ProcessUsageEvent processes a usage event synchronously.
func (s *WorkflowServer) ProcessUsageEvent(ctx context.Context, req *workflowv1.ProcessUsageEventRequest) (*workflowv1.ProcessUsageEventResponse, error) {
	receivedAt := time.Now().UTC()
	result, err := s.svc.ProcessUsageEvent(ctx, req.GetEventId(), req.GetOrgId(), req.GetEvent(), receivedAt)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &workflowv1.ProcessUsageEventResponse{
		Status:  result.Status,
		CostUsd: result.CostUSD,
	}, nil
}