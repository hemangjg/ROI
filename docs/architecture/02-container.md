# C4 — Container Diagram

Deployable containers in the Phase 1 local stack. Production adds Cloudflare edge and EKS; see [10-infrastructure.md](../phase-0.5/10-infrastructure.md).

```mermaid
flowchart TB
    subgraph clients [Clients]
        SDKs["SDKs TypeScript / Python / Go"]
        Dashboard["Next.js Dashboard"]
    end

    subgraph edge [Edge — production]
        CF["Cloudflare"]
    end

    Envoy["Envoy Gateway"]

    subgraph authLayer [Identity and Authorization]
        Identity["identity-service"]
        OpenFGA["OpenFGA"]
    end

    subgraph services [Application Services]
        Ingestion["ingestion-service"]
        Management["management-service"]
        Analytics["analytics-service"]
        Pricing["pricing-service"]
        Workflow["workflow-service"]
    end

    subgraph messaging [Messaging and Workflows]
        Kafka["Kafka"]
        Temporal["Temporal"]
    end

    subgraph data [Data Stores]
        PG["PostgreSQL"]
        CH["ClickHouse"]
        Redis["Redis"]
    end

    subgraph obs [Observability]
        OTel["OTel Collector"]
        Prom["Prometheus"]
        Grafana["Grafana"]
        Jaeger["Jaeger"]
        Loki["Loki"]
    end

    SDKs --> CF
    Dashboard --> CF
    CF --> Envoy
    SDKs -.->|"local dev"| Envoy
    Dashboard -.->|"local dev"| Envoy

    Envoy --> Identity
    Envoy --> Ingestion
    Envoy --> Management
    Envoy --> Analytics
    Envoy --> Dashboard

    Identity --> OpenFGA
    Identity --> PG
    Identity --> Redis
    Ingestion --> PG
    Ingestion --> Kafka
    Ingestion --> Redis
    Management --> PG
    Management --> OpenFGA
    Analytics --> CH
    Analytics --> Redis
    Analytics --> OpenFGA
    Pricing --> PG
    Workflow --> Temporal
    Workflow --> Kafka
    Workflow --> Pricing
    Workflow --> CH
    Workflow --> PG

    services --> OTel
    Identity --> OTel
    OTel --> Prom
    OTel --> Jaeger
    OTel --> Loki
    Prom --> Grafana
```

## Local ports (Docker Compose)

| Container | Port | Notes |
|-----------|------|-------|
| identity-service | 8080 | `/v1/auth` via Envoy |
| ingestion-service | 8081 | `/v1/events` via Envoy |
| management-service | 8082 | `/v1/orgs` via Envoy |
| analytics-service | 8083 | `/v1/orgs/*/spend` via Envoy |
| pricing-service | 8084 | gRPC only (internal) |
| workflow-service | 8085 | Temporal worker |
| web (dashboard) | 3000 | Next.js |
| Envoy HTTP | 8888 | API gateway |
| Envoy HTTPS API | 8443 | `api.ai-finops.local` |
| Envoy HTTPS App | 8444 | `app.ai-finops.local` |

## Communication summary

| From | To | Protocol | Phase |
|------|-----|----------|-------|
| SDK / Dashboard | Services | REST via Envoy | 1 (routing); 2 (handlers) |
| Envoy | identity-service | ext_authz gRPC | 2 |
| Ingestion | Kafka | Kafka | 2 |
| Workflow | Pricing | gRPC | 2 |
| All services | OTel Collector | OTLP | 1 |