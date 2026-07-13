#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CLICKHOUSE_URL="${CLICKHOUSE_URL:-http://localhost:8123}"
DB="${CLICKHOUSE_DB:-finops}"

apply_file() {
  local file="$1"
  echo "    applying $(basename "$file")"
  local statement=""
  while IFS= read -r line || [[ -n "$line" ]]; do
    [[ "$line" =~ ^[[:space:]]*-- ]] && continue
    [[ -z "${line// }" ]] && continue
    statement+="${line}"$'\n'
    if [[ "$line" == *";" ]]; then
      if ! curl -fsS "${CLICKHOUSE_URL}/?database=${DB}" --data-binary "$statement"; then
        echo "ERROR: failed statement from $(basename "$file"):"
        echo "$statement"
        exit 1
      fi
      statement=""
    fi
  done < "$file"
}

if ! curl -fsS "${CLICKHOUSE_URL}/ping" >/dev/null 2>&1; then
  echo "ClickHouse not reachable at ${CLICKHOUSE_URL} — skipping migrations"
  exit 0
fi

curl -fsS "${CLICKHOUSE_URL}/" --data-binary "CREATE DATABASE IF NOT EXISTS ${DB}" >/dev/null

for file in "$ROOT"/db/clickhouse/migrations/*.sql; do
  apply_file "$file"
done

echo "ClickHouse migrations applied."