package grpcapi

import (
	"context"

	analyticsapp "github.com/ai-finops/ai-finops/services/analytics-service/internal/application/analytics"
	analyticsv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/analytics/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AnalyticsServer implements analytics.v1.AnalyticsService.
type AnalyticsServer struct {
	analyticsv1.UnimplementedAnalyticsServiceServer
	svc *analyticsapp.Service
}

// NewAnalyticsServer constructs the gRPC analytics service.
func NewAnalyticsServer(svc *analyticsapp.Service) *AnalyticsServer {
	return &AnalyticsServer{svc: svc}
}

// GetSpendSummary returns organization spend summary metrics.
func (s *AnalyticsServer) GetSpendSummary(ctx context.Context, req *analyticsv1.GetSpendSummaryRequest) (*analyticsv1.GetSpendSummaryResponse, error) {
	params, err := analyticsapp.ParseQueryParams(
		req.GetOrgId(),
		timestampToString(req.GetFrom()),
		timestampToString(req.GetTo()),
		req.GetTeamId(),
	)
	if err != nil {
		return nil, toGRPCError(err)
	}

	result, err := s.svc.GetSpendSummary(ctx, params)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &analyticsv1.GetSpendSummaryResponse{
		TodayUsd:          result.TodayUSD,
		MtdUsd:            result.MTDUSD,
		PreviousPeriodUsd: result.PreviousPeriodUSD,
		ChangePct:         result.ChangePct,
	}, nil
}

// GetSpendTimeSeries returns daily spend points for the requested range.
func (s *AnalyticsServer) GetSpendTimeSeries(ctx context.Context, req *analyticsv1.GetSpendTimeSeriesRequest) (*analyticsv1.GetSpendTimeSeriesResponse, error) {
	params, err := analyticsapp.ParseQueryParams(
		req.GetOrgId(),
		timestampToString(req.GetFrom()),
		timestampToString(req.GetTo()),
		req.GetTeamId(),
	)
	if err != nil {
		return nil, toGRPCError(err)
	}

	result, err := s.svc.GetSpendTimeSeries(ctx, params, req.GetGranularity())
	if err != nil {
		return nil, toGRPCError(err)
	}

	points := make([]*analyticsv1.TimeSeriesPoint, 0, len(result.Data))
	for _, point := range result.Data {
		points = append(points, &analyticsv1.TimeSeriesPoint{
			Date:    point.Date,
			CostUsd: point.CostUSD,
		})
	}
	return &analyticsv1.GetSpendTimeSeriesResponse{Data: points}, nil
}

func timestampToString(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().UTC().Format("2006-01-02T15:04:05Z07:00")
}