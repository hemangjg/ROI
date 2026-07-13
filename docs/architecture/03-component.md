# C4 — Component Diagram

Phase 1 monorepo structure and the internal layout shared by every Go service shell.

## Monorepo components

```mermaid
flowchart TB
    subgraph apps [apps/]
        Web["web — Next.js dashboard"]
    end

    subgraph services [services/]
        Identity["identity-service"]
        Ingestion["ingestion-service"]
        Management["management-service"]
        Analytics["analytics-service"]
        Pricing["pricing-service"]
        Workflow["workflow-service"]
    end

    subgraph packages [packages/]
        ProtoPkg["proto — generated stubs"]
        SDKCore["sdk-core — HTTP client"]
        SDKTS["sdk-typescript — public SDK"]
        APISchemas["api-schemas — OpenAPI types"]
        SharedGo["auth, config, logger, shared, validation, types, events"]
    end

    subgraph contracts [contracts/]
        BufProto["proto — Buf module"]
        OpenAPI["api/openapi/v1.yaml"]
    end

    subgraph dataLayer [db/]
        Atlas["atlas — PostgreSQL migrations"]
        SQLC["sql — SQLC queries"]
    end

    subgraph tooling [tooling]
        Nx["Nx + pnpm"]
        CI["GitHub Actions"]
        Compose["Docker Compose"]
        EnvoyCfg["infra/envoy"]
    end

    Web --> APISchemas
    Web --> SDKTS
    SDKTS --> SDKCore
    SDKTS --> APISchemas
    services --> ProtoPkg
    services --> SharedGo
    BufProto -->|"buf generate"| ProtoPkg
    Atlas --> SQLC
    services --> SQLC
    Nx --> apps
    Nx --> services
    Nx --> packages
    Compose --> services
    EnvoyCfg --> Compose
```

## Go service internal components

Each service under `services/{name}/` follows the same hexagonal layout (Phase 1 shells wire health + readiness only):

```mermaid
flowchart TB
    subgraph interfaces [interfaces/]
        HTTP["http — Chi routes, middleware"]
        GRPC["grpc — gRPC server stubs Phase 2"]
    end

    subgraph application [application/]
        UseCases["use cases — Phase 2"]
    end

    subgraph domain [domain/]
        Entities["entities + repository ports — Phase 2"]
    end

    subgraph infrastructure [infrastructure/]
        Postgres["postgres — SQLC repos"]
        Redis["redis — cache Phase 2"]
        Config["config — Koanf"]
        OTel["otel — tracing middleware"]
    end

    HTTP --> UseCases
    GRPC --> UseCases
    UseCases --> Entities
    UseCases --> Postgres
    UseCases --> Redis
    HTTP --> OTel
```

## Nx orchestration

| Target | Projects | Purpose |
|--------|----------|---------|
| `build` | services, packages | Compile Go binaries and TS packages |
| `test` | services, packages, web | Unit and integration tests |
| `serve` | services, web | Local dev with hot reload |
| `proto-generate` | workspace | `buf lint` + `buf generate` |
| `sqlc-generate` | workspace | Regenerate SQLC from queries |
| `dev` | workspace | `scripts/dev.sh` — preflight + concurrent serve |