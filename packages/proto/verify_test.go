package proto_test

import (
	"testing"

	analyticsv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/analytics/v1"
	authv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/auth/v1"
	healthv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/bootstrap/health/v1"
	ingestionv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/ingestion/v1"
	managementv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/management/v1"
	pricingv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/pricing/v1"
	workflowv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/workflow/v1"
	"google.golang.org/protobuf/proto"
)

func assertRoundTrip(t *testing.T, msg proto.Message) {
	t.Helper()

	b, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal %T: %v", msg, err)
	}

	out := proto.Clone(msg)
	out.ProtoReflect().SetUnknown(nil)
	if err := proto.Unmarshal(b, out); err != nil {
		t.Fatalf("unmarshal %T: %v", msg, err)
	}

	if !proto.Equal(msg, out) {
		t.Fatalf("round-trip mismatch for %T", msg)
	}
}

func TestHealthCheckResponseRoundTrip(t *testing.T) {
	assertRoundTrip(t, &healthv1.HealthCheckResponse{Status: "ok"})
}

func TestAuthValidateTokenResponseRoundTrip(t *testing.T) {
	assertRoundTrip(t, &authv1.ValidateTokenResponse{
		Valid:  true,
		UserId: "user-1",
		OrgId:  "org-1",
		Roles:  []string{"admin"},
	})
}

func TestIngestionUsageEventRoundTrip(t *testing.T) {
	assertRoundTrip(t, &ingestionv1.UsageEvent{
		IdempotencyKey: "key-1",
		Provider:       "openai",
		Model:          "gpt-4o",
		InputTokens:    100,
		OutputTokens:   50,
		Metadata: &ingestionv1.UsageEventMetadata{
			TeamId:  "team-1",
			Project: "billing",
		},
	})
}

func TestPricingCalculateCostResponseRoundTrip(t *testing.T) {
	assertRoundTrip(t, &pricingv1.CalculateCostResponse{
		CostUsd:  "0.0125",
		IsPriced: true,
		ModelId:  "model-1",
	})
}

func TestAnalyticsSpendSummaryResponseRoundTrip(t *testing.T) {
	assertRoundTrip(t, &analyticsv1.GetSpendSummaryResponse{
		TodayUsd:          "12.50",
		MtdUsd:            "340.00",
		PreviousPeriodUsd: "300.00",
		ChangePct:         13.3,
	})
}

func TestManagementCreateApiKeyResponseRoundTrip(t *testing.T) {
	assertRoundTrip(t, &managementv1.CreateApiKeyResponse{
		KeyId:     "key-1",
		ApiKey:    "sk_live_xxx",
		KeyPrefix: "sk_live",
	})
}

func TestWorkflowProcessUsageEventResponseRoundTrip(t *testing.T) {
	assertRoundTrip(t, &workflowv1.ProcessUsageEventResponse{
		Status:  "processed",
		CostUsd: "0.0125",
	})
}