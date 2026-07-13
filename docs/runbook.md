# Runbook — local engineering base

## Prerequisites

- Docker Desktop running
- Go 1.25+, Node 22+, pnpm 9+, buf, atlas, sqlc, golangci-lint

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
export PATH="$(go env GOPATH)/bin:$PATH"
```

## Bring-up

```bash
cd "/Users/hemangjg/Hemangs stuff/Projects/ai-finops"
pnpm install
go work sync
docker compose --profile services up -d --build
bash scripts/init-local.sh   # migrations + kafka topic + optional workflow bounce
bash scripts/e2e-smoke.sh
```

Dashboard (optional):

```bash
cd apps/web && NEXT_PUBLIC_USE_MOCK_API=false pnpm dev
# http://localhost:3000
```

API base: `http://localhost:8888/v1`

## Quality gate

```bash
bash scripts/quality-gate.sh           # offline
WITH_E2E=1 bash scripts/quality-gate.sh  # offline + e2e
```

## Failure checklist

| Symptom | Fix |
|---------|-----|
| Prettier fail | `pnpm format` |
| golangci errcheck | fix ignored errors |
| Postgres not reachable | `docker compose --profile core up -d` |
| E2E spend stays 0 | `docker compose restart workflow-service` then re-run e2e |
| Spend 404 via Envoy | check envoy spend regex is full-path (`.../spend(/.*)?$`) |
| Kafka unknown topic | `bash scripts/init-local.sh` or wait for `EnsureTopic` on service start |

## Verify services

```bash
for p in 8080 8081 8082 8083 8084 8085; do curl -fsS "http://localhost:$p/healthz"; echo " :$p"; done
curl -fsS http://localhost:8888/gateway/healthz
```

## Status log

See [status.md](./status.md) for last verified pass and engineering sign-off.
