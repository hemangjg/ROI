/**
 * @ai-finops/proto — protobuf types for AI FinOps domain contracts.
 * Run `pnpm proto:generate` to regenerate packages/proto/gen/{go,ts}.
 */
export type {
  HealthCheckRequest,
  HealthCheckResponse,
} from "../gen/ts/bootstrap/health/v1/health_pb.js";

export {
  HealthCheckRequestSchema,
  HealthCheckResponseSchema,
} from "../gen/ts/bootstrap/health/v1/health_pb.js";

export type {
  ValidateTokenRequest,
  ValidateTokenResponse,
  ValidateApiKeyRequest,
  ValidateApiKeyResponse,
} from "../gen/ts/auth/v1/auth_pb.js";

export {
  ValidateTokenRequestSchema,
  ValidateTokenResponseSchema,
  ValidateApiKeyRequestSchema,
  ValidateApiKeyResponseSchema,
  AuthService,
} from "../gen/ts/auth/v1/auth_pb.js";

export type {
  UsageEvent,
  UsageEventMetadata,
  IngestEventRequest,
  IngestEventResponse,
  UsageEventCreated,
} from "../gen/ts/ingestion/v1/event_pb.js";

export {
  UsageEventSchema,
  UsageEventMetadataSchema,
  IngestEventRequestSchema,
  IngestEventResponseSchema,
  UsageEventCreatedSchema,
  IngestionService,
} from "../gen/ts/ingestion/v1/event_pb.js";

export type {
  CalculateCostRequest,
  CalculateCostResponse,
  PriceEntry,
  SyncPricingCatalogRequest,
  SyncPricingCatalogResponse,
} from "../gen/ts/pricing/v1/pricing_pb.js";

export {
  CalculateCostRequestSchema,
  CalculateCostResponseSchema,
  PriceEntrySchema,
  SyncPricingCatalogRequestSchema,
  SyncPricingCatalogResponseSchema,
  PricingService,
} from "../gen/ts/pricing/v1/pricing_pb.js";

export type {
  GetSpendSummaryRequest,
  GetSpendSummaryResponse,
  TimeSeriesPoint,
  GetSpendTimeSeriesRequest,
  GetSpendTimeSeriesResponse,
} from "../gen/ts/analytics/v1/spend_pb.js";

export {
  GetSpendSummaryRequestSchema,
  GetSpendSummaryResponseSchema,
  TimeSeriesPointSchema,
  GetSpendTimeSeriesRequestSchema,
  GetSpendTimeSeriesResponseSchema,
  AnalyticsService,
} from "../gen/ts/analytics/v1/spend_pb.js";

export type {
  CreateOrgRequest,
  CreateOrgResponse,
  CreateApiKeyRequest,
  CreateApiKeyResponse,
} from "../gen/ts/management/v1/org_pb.js";

export {
  CreateOrgRequestSchema,
  CreateOrgResponseSchema,
  CreateApiKeyRequestSchema,
  CreateApiKeyResponseSchema,
  ManagementService,
} from "../gen/ts/management/v1/org_pb.js";

export type {
  ProcessUsageEventRequest,
  ProcessUsageEventResponse,
} from "../gen/ts/workflow/v1/workflow_pb.js";

export {
  ProcessUsageEventRequestSchema,
  ProcessUsageEventResponseSchema,
  WorkflowService,
} from "../gen/ts/workflow/v1/workflow_pb.js";
