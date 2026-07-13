#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

export DATABASE_URL="${DATABASE_URL:-postgres://finops:finops@localhost:5432/finops?sslmode=disable}"
export ATLAS_DEV_URL="${ATLAS_DEV_URL:-postgres://finops:finops@localhost:5432/atlas_dev?sslmode=disable}"

wait_for_postgres() {
  echo "==> waiting for postgres at localhost:5432"
  for _ in $(seq 1 60); do
    if command -v pg_isready >/dev/null 2>&1 && pg_isready -d "$DATABASE_URL" >/dev/null 2>&1; then
      return 0
    fi
    if command -v nc >/dev/null 2>&1 && nc -z localhost 5432 >/dev/null 2>&1; then
      if docker compose -f "$ROOT/docker-compose.yml" exec -T postgres pg_isready -U finops -d finops >/dev/null 2>&1; then
        return 0
      fi
    fi
    sleep 2
  done
  echo "postgres not ready after 120s"
  return 1
}

wait_for_clickhouse() {
  echo "==> waiting for clickhouse at localhost:8123"
  for _ in $(seq 1 60); do
    if curl -fsS "http://localhost:8123/ping" >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done
  echo "clickhouse not ready after 120s"
  return 1
}

if ! command -v atlas >/dev/null 2>&1; then
  echo "atlas not installed; skipping postgres migrations"
else
  wait_for_postgres
  cd "$ROOT/db/atlas"
  echo "==> atlas migrate apply"
  atlas migrate apply --env local
fi

if curl -fsS "http://localhost:8123/ping" >/dev/null 2>&1; then
  echo "==> clickhouse migrations"
  bash "$ROOT/scripts/apply-clickhouse-migrations.sh"
else
  echo "==> skipping clickhouse migrations (not reachable)"
fi

# Ensure Kafka topic exists so workflow consumer gets partitions on first start.
if docker compose -f "$ROOT/docker-compose.yml" ps --status running 2>/dev/null | grep -q kafka; then
  echo "==> ensuring kafka topic usage.events.raw"
  docker compose -f "$ROOT/docker-compose.yml" exec -T kafka \
    /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 \
    --create --if-not-exists --topic usage.events.raw \
    --partitions 1 --replication-factor 1 >/dev/null || true
  # kafka-go can end up with 0 assigned partitions after broker restarts; bounce consumer.
  if docker compose -f "$ROOT/docker-compose.yml" ps --status running 2>/dev/null | grep -q workflow-service; then
    echo "==> restarting workflow-service to rejoin kafka consumer group"
    docker compose -f "$ROOT/docker-compose.yml" restart workflow-service >/dev/null || true
  fi
fi

echo "Local initialization complete."
echo "Pricing seed data is added in a later Phase 1b workstream."