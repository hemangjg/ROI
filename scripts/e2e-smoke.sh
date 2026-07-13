#!/usr/bin/env bash
# End-to-end Phase 1b smoke: register → API key → ingest → process → spend summary.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

BASE_URL="${BASE_URL:-http://localhost:8888/v1}"
TIMEOUT_SECS="${E2E_TIMEOUT_SECS:-90}"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "ERROR: required command not found: $1" >&2
    exit 1
  }
}

need_cmd curl
need_cmd python3

json_field() {
  local json="$1"
  local field="$2"
  python3 -c 'import json,sys; d=json.loads(sys.argv[1]); print(d.get(sys.argv[2],""))' "$json" "$field"
}

http_json() {
  # args: method url [curl args...]
  local method="$1"
  local url="$2"
  shift 2
  local tmp
  tmp="$(mktemp)"
  local code
  code="$(curl -sS -o "$tmp" -w "%{http_code}" -X "$method" "$url" "$@")"
  BODY="$(cat "$tmp")"
  rm -f "$tmp"
  HTTP_CODE="$code"
}

echo "==> e2e smoke against ${BASE_URL}"

echo "==> health via Envoy gateway"
http_json GET "http://localhost:8888/gateway/healthz" || true
if [[ "${HTTP_CODE:-}" != "200" ]]; then
  # some builds expose health only on service ports
  curl -fsS "http://localhost:8080/healthz" >/dev/null
fi

EMAIL="e2e-$(python3 -c 'import uuid; print(uuid.uuid4())')@example.com"
PASSWORD="e2e-password-change-me"
ORG="E2E Org $(python3 -c 'import uuid; print(uuid.uuid4().hex[:8])')"

echo "==> register ${EMAIL}"
http_json POST "${BASE_URL}/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\",\"name\":\"E2E User\",\"org_name\":\"${ORG}\"}"
if [[ "$HTTP_CODE" != "201" ]]; then
  echo "register failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi

ACCESS_TOKEN="$(json_field "$BODY" access_token)"
ORG_ID="$(json_field "$BODY" org_id)"
if [[ -z "$ACCESS_TOKEN" || -z "$ORG_ID" ]]; then
  echo "register response missing tokens/org_id: ${BODY}" >&2
  exit 1
fi
echo "    org_id=${ORG_ID}"

echo "==> create API key"
http_json POST "${BASE_URL}/orgs/${ORG_ID}/api-keys/" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H 'Content-Type: application/json' \
  -d '{"name":"e2e-ingest-key"}'
if [[ "$HTTP_CODE" != "201" ]]; then
  echo "create api key failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi

API_KEY="$(json_field "$BODY" api_key)"
if [[ -z "$API_KEY" ]]; then
  echo "api key missing in response: ${BODY}" >&2
  exit 1
fi
echo "    api key created (prefix hidden)"

IDEM="e2e-$(python3 -c 'import uuid; print(uuid.uuid4())')"
OCCURRED_AT="$(python3 -c 'from datetime import datetime, timezone; print(datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"))')"

echo "==> ingest usage event"
http_json POST "${BASE_URL}/events" \
  -H "Authorization: Bearer ${API_KEY}" \
  -H 'Content-Type: application/json' \
  -d "{\"idempotency_key\":\"${IDEM}\",\"provider\":\"openai\",\"model\":\"gpt-4o\",\"input_tokens\":1200,\"output_tokens\":340,\"occurred_at\":\"${OCCURRED_AT}\"}"
if [[ "$HTTP_CODE" != "202" ]]; then
  echo "ingest failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi

EVENT_ID="$(json_field "$BODY" event_id)"
STATUS="$(json_field "$BODY" status)"
echo "    event_id=${EVENT_ID} status=${STATUS}"
if [[ -z "$EVENT_ID" ]]; then
  echo "missing event_id" >&2
  exit 1
fi

echo "==> wait for outbox/workflow processing (timeout ${TIMEOUT_SECS}s)"
deadline=$((SECONDS + TIMEOUT_SECS))
PROCESSED=0
while (( SECONDS < deadline )); do
  http_json GET "${BASE_URL}/events/${EVENT_ID}" \
    -H "Authorization: Bearer ${API_KEY}" || true
  if [[ "${HTTP_CODE:-}" == "200" ]]; then
    STATUS="$(json_field "$BODY" status)"
    COST="$(json_field "$BODY" cost_usd)"
    if [[ "$STATUS" == "processed" || "$STATUS" == "priced" || -n "$COST" ]]; then
      PROCESSED=1
      echo "    event status=${STATUS} cost_usd=${COST:-n/a}"
      break
    fi
  fi
  # also poll spend summary once pipeline may have written ClickHouse
  FROM="$(python3 -c 'from datetime import date,timedelta; print((date.today()-timedelta(days=1)).isoformat())')"
  TO="$(python3 -c 'from datetime import date,timedelta; print((date.today()+timedelta(days=1)).isoformat())')"
  http_json GET "${BASE_URL}/orgs/${ORG_ID}/spend/summary?from=${FROM}&to=${TO}" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" || true
  if [[ "${HTTP_CODE:-}" == "200" ]]; then
    TOTAL="$(python3 -c 'import json,sys
d=json.loads(sys.argv[1])
for k in ("mtd_usd","today_usd","total_cost_usd","total","cost_usd"):
  v=d.get(k)
  if v not in (None, "", "0", "0.0", "0.00", 0):
    print(v); break
' "$BODY")"
    if [[ -n "$TOTAL" ]]; then
      PROCESSED=1
      echo "    spend summary total=${TOTAL}"
      break
    fi
  fi
  sleep 2
done

if [[ "$PROCESSED" != "1" ]]; then
  echo "WARN: event did not show non-zero spend within ${TIMEOUT_SECS}s"
  echo "      last event body: ${BODY:-}"
  # still fail hard — e2e must prove the loop
  exit 1
fi

FROM="$(python3 -c 'from datetime import date,timedelta; print((date.today()-timedelta(days=1)).isoformat())')"
TO="$(python3 -c 'from datetime import date,timedelta; print((date.today()+timedelta(days=1)).isoformat())')"
echo "==> spend summary"
http_json GET "${BASE_URL}/orgs/${ORG_ID}/spend/summary?from=${FROM}&to=${TO}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
if [[ "$HTTP_CODE" != "200" ]]; then
  echo "spend summary failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi
echo "    ${BODY}"

echo ""
echo "E2E smoke PASSED"
echo "  register → api key → ingest → process/spend verified"
