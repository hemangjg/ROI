package pricing

import (
	"context"
	"fmt"
	"time"

	pricingv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/pricing/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CostQuote is a pricing-service cost calculation result.
type CostQuote struct {
	CostUSD  string
	IsPriced bool
	ModelID  string
}

// Client wraps pricing-service gRPC calls.
type Client struct {
	conn   *grpc.ClientConn
	client pricingv1.PricingServiceClient
}

// NewClient dials the pricing-service gRPC endpoint.
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial pricing grpc: %w", err)
	}
	return &Client{
		conn:   conn,
		client: pricingv1.NewPricingServiceClient(conn),
	}, nil
}

// CalculateCost requests priced cost for a usage event.
func (c *Client) CalculateCost(
	ctx context.Context,
	provider, model string,
	inputTokens, outputTokens, cacheReadTokens, cacheWriteTokens uint32,
	occurredAt time.Time,
) (CostQuote, error) {
	resp, err := c.client.CalculateCost(ctx, &pricingv1.CalculateCostRequest{
		Provider:         provider,
		Model:            model,
		InputTokens:      inputTokens,
		OutputTokens:     outputTokens,
		CacheReadTokens:  cacheReadTokens,
		CacheWriteTokens: cacheWriteTokens,
		OccurredAt:       timestamppb.New(occurredAt.UTC()),
	})
	if err != nil {
		return CostQuote{}, fmt.Errorf("calculate cost rpc: %w", err)
	}

	return CostQuote{
		CostUSD:  resp.GetCostUsd(),
		IsPriced: resp.GetIsPriced(),
		ModelID:  resp.GetModelId(),
	}, nil
}

// Close closes the gRPC connection.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}