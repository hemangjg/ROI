# C4 — Context Diagram

Shows external actors and systems that interact with AI FinOps.

```mermaid
flowchart TB
    subgraph external [External Actors]
        Engineer["Platform Engineer"]
        FinOps["FinOps Manager"]
        CISO["CISO / Security"]
    end

    subgraph customer [Customer Environment]
        Apps["Customer Applications"]
        SDK["AI FinOps SDKs"]
    end

    subgraph platform [AI FinOps Platform]
        System["AI FinOps System"]
    end

    subgraph providers [LLM Providers]
        OpenAI["OpenAI"]
        Anthropic["Anthropic"]
        Others["Other Providers"]
    end

    Engineer -->|"Installs SDK, configures keys"| SDK
    FinOps -->|"Views spend dashboards"| System
    CISO -->|"Reviews audit and policies"| System
    Apps -->|"LLM API calls"| OpenAI
    Apps -->|"LLM API calls"| Anthropic
    Apps -->|"LLM API calls"| Others
    SDK -->|"Usage events via REST"| System
    Apps -->|"Optional direct ingest"| System
```

## Responsibilities

| Actor | Goal |
|-------|------|
| Platform Engineer | Instrument apps with SDKs, route traffic through gateway |
| FinOps Manager | Track spend by team, model, and provider |
| CISO | Enforce governance, audit access, review policies |
| Customer Apps | Call LLM providers; emit usage telemetry |
| LLM Providers | Billable inference APIs (external to AI FinOps) |