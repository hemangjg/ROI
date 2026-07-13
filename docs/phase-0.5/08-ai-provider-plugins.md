# 08 — AI Provider Plugins

## Plugin Architecture

Provider plugins normalize LLM usage data and pricing across vendors. They live in the **Pricing Service** and are invoked by the **Workflow Service** during cost calculation and catalog sync.

---

## Interface (Go)

```go
// services/pricing-service/internal/domain/plugin.go

type ProviderPlugin interface {
    Name() string
    Normalize(raw RawUsageEvent) (NormalizedEvent, error)
    SupportedModels() []ModelDescriptor
    FetchPricing(ctx context.Context) ([]PriceEntry, error)
}

type RawUsageEvent struct {
    Provider    string
    Model       string
    RawMetadata map[string]any
}

type NormalizedEvent struct {
    Provider         string
    Model            string
    InputTokens      uint32
    OutputTokens     uint32
    CacheReadTokens  uint32
    CacheWriteTokens uint32
}

type ModelDescriptor struct {
    Name        string
    DisplayName string
    Provider    string
}

type PriceEntry struct {
    ModelID              string
    InputPricePer1M      decimal.Decimal
    OutputPricePer1M     decimal.Decimal
    CacheReadPricePer1M  decimal.Decimal
    CacheWritePricePer1M decimal.Decimal
    EffectiveFrom        time.Time
}
```

Plugins register via Wire DI at startup:

```go
var PluginSet = wire.NewSet(
    openai.NewPlugin,
    anthropic.NewPlugin,
    plugin.NewRegistry,
)
```

---

## Provider Matrix

| Provider | Phase | Plugin package | SDK normalization |
|----------|-------|---------------|-------------------|
| OpenAI | 1 | `plugin/openai` | `response.usage` fields |
| Anthropic | 1 | `plugin/anthropic` | `usage.input_tokens`, `output_tokens` |
| Google Gemini | 2 | `plugin/gemini` | `usageMetadata` |
| Grok (xAI) | 2 | `plugin/grok` | OpenAI-compatible usage |
| DeepSeek | 2 | `plugin/deepseek` | OpenAI-compatible usage |
| Mistral | 2 | `plugin/mistral` | `usage` object |
| Azure OpenAI | 3 | `plugin/azureopenai` | Deployment name mapping |
| AWS Bedrock | 3 | `plugin/bedrock` | Model ID + token counts |
| Vertex AI | 3 | `plugin/vertexai` | `usageMetadata` |
| Ollama | 3 | `plugin/ollama` | Zero cost; token counts only |

---

## Phase 1 Plugins

### OpenAI

Models: `gpt-4o`, `gpt-4o-mini`, `gpt-4-turbo`, `o1`, `o1-mini`

Pricing source: OpenAI pricing page (weekly sync via `SyncPricingCatalog` workflow).

Token fields: `prompt_tokens`, `completion_tokens`, `prompt_tokens_details.cached_tokens`

### Anthropic

Models: `claude-3-5-sonnet-20241022`, `claude-3-5-haiku-20241022`, `claude-3-opus-20240229`

Pricing source: Anthropic pricing page.

Token fields: `usage.input_tokens`, `usage.output_tokens`, `usage.cache_read_input_tokens`, `usage.cache_creation_input_tokens`

---

## SDK Instrumentation Pattern

SDK wraps provider client; emits normalized event after response:

```typescript
// packages/sdk-typescript/src/openai.ts
const response = await openai.chat.completions.create(params);
await finops.track({
  idempotency_key: requestId,
  provider: 'openai',
  model: params.model,
  input_tokens: response.usage.prompt_tokens,
  output_tokens: response.usage.completion_tokens,
  cache_read_tokens: response.usage.prompt_tokens_details?.cached_tokens ?? 0,
  occurred_at: new Date().toISOString(),
  metadata: { team_id, project, feature },
});
```

SDK never blocks the LLM call — `track()` is fire-and-forget with local retry buffer.

---

## Unknown Model Handling

1. Pricing Service returns `is_priced: false`
2. Workflow writes event with `cost_usd: 0`, `is_priced: 0`
3. Metric `unpriced_events_total{provider, model}` incremented
4. Management Service shows admin alert when unpriced events > 0 in 24h

---

## Future: Custom Plugins (Phase 4)

Enterprise customers can register custom provider plugins via API:
- Upload pricing JSON
- Map custom model names
- Stored in PostgreSQL `custom_providers` table