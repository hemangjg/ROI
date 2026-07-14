#!/usr/bin/env bash
# Phase 2 budgets smoke: register → team → budget → evaluate soft alert → list/ack
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/opt/homebrew/bin:${PATH:-}"

BASE_URL="${BASE_URL:-http://localhost:8888/v1}"

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

echo "==> e2e budgets against ${BASE_URL}"

EMAIL="budget-e2e-$(python3 -c 'import uuid; print(uuid.uuid4())')@example.com"
PASSWORD="e2e-password-change-me"
ORG="Budget E2E $(python3 -c 'import uuid; print(uuid.uuid4().hex[:8])')"

echo "==> register"
http_json POST "${BASE_URL}/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\",\"name\":\"Budget E2E\",\"org_name\":\"${ORG}\"}"
if [[ "$HTTP_CODE" != "201" ]]; then
  echo "register failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi

ACCESS_TOKEN="$(json_field "$BODY" access_token)"
ORG_ID="$(json_field "$BODY" org_id)"
AUTH=(-H "Authorization: Bearer ${ACCESS_TOKEN}" -H 'Content-Type: application/json')
echo "    org_id=${ORG_ID}"

echo "==> create team"
http_json POST "${BASE_URL}/orgs/${ORG_ID}/teams/" \
  "${AUTH[@]}" \
  -d '{"name":"Platform"}'
if [[ "$HTTP_CODE" != "201" ]]; then
  echo "create team failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi
TEAM_ID="$(json_field "$BODY" id)"
if [[ -z "$TEAM_ID" ]]; then
  echo "team id missing: ${BODY}" >&2
  exit 1
fi
echo "    team_id=${TEAM_ID}"

PERIOD="$(python3 -c 'from datetime import date; d=date.today(); print(f"{d.year}-{d.month:02d}")')"

echo "==> upsert budget \$1.00 soft 50% period ${PERIOD}"
http_json PUT "${BASE_URL}/orgs/${ORG_ID}/budgets" \
  "${AUTH[@]}" \
  -d "{\"team_id\":\"${TEAM_ID}\",\"amount_usd\":\"1.00\",\"soft_threshold_pct\":50,\"period_start\":\"${PERIOD}\"}"
if [[ "$HTTP_CODE" != "200" && "$HTTP_CODE" != "201" ]]; then
  echo "upsert budget failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi
BUDGET_ID="$(json_field "$BODY" id)"
echo "    budget_id=${BUDGET_ID}"

echo "==> evaluate with spend 0.80 (over soft 50%)"
http_json POST "${BASE_URL}/orgs/${ORG_ID}/budgets/evaluate" \
  "${AUTH[@]}" \
  -d "{\"team_id\":\"${TEAM_ID}\",\"spend_usd\":\"0.80\",\"period_start\":\"${PERIOD}\"}"
if [[ "$HTTP_CODE" != "200" ]]; then
  echo "evaluate failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi
OVER="$(json_field "$BODY" over_soft_threshold)"
FIRED="$(json_field "$BODY" alert_fired)"
echo "    over_soft_threshold=${OVER} alert_fired=${FIRED}"
if [[ "$OVER" != "True" && "$OVER" != "true" ]]; then
  echo "expected soft threshold breach: ${BODY}" >&2
  exit 1
fi
if [[ "$FIRED" != "True" && "$FIRED" != "true" ]]; then
  echo "expected alert_fired: ${BODY}" >&2
  exit 1
fi

echo "==> list open alerts"
http_json GET "${BASE_URL}/orgs/${ORG_ID}/budget-alerts?status=open" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
if [[ "$HTTP_CODE" != "200" ]]; then
  echo "list alerts failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi
ALERT_COUNT="$(python3 -c 'import json,sys; print(len(json.loads(sys.argv[1])))' "$BODY")"
if [[ "$ALERT_COUNT" -lt 1 ]]; then
  echo "expected at least one open alert: ${BODY}" >&2
  exit 1
fi
ALERT_ID="$(python3 -c 'import json,sys; print(json.loads(sys.argv[1])[0]["id"])' "$BODY")"
echo "    open_alerts=${ALERT_COUNT} first=${ALERT_ID}"

echo "==> acknowledge alert"
http_json POST "${BASE_URL}/orgs/${ORG_ID}/budget-alerts/${ALERT_ID}/ack" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
if [[ "$HTTP_CODE" != "200" ]]; then
  echo "ack failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi
STATUS="$(json_field "$BODY" status)"
if [[ "$STATUS" != "acknowledged" ]]; then
  echo "expected status acknowledged: ${BODY}" >&2
  exit 1
fi

echo "==> list budgets"
http_json GET "${BASE_URL}/orgs/${ORG_ID}/budgets?period=${PERIOD}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
if [[ "$HTTP_CODE" != "200" ]]; then
  echo "list budgets failed HTTP ${HTTP_CODE}: ${BODY}" >&2
  exit 1
fi
COUNT="$(python3 -c 'import json,sys; print(len(json.loads(sys.argv[1])))' "$BODY")"
if [[ "$COUNT" -lt 1 ]]; then
  echo "expected budget listed: ${BODY}" >&2
  exit 1
fi

echo ""
echo "E2E budgets PASSED"
echo "  register → team → budget → soft alert → ack"
