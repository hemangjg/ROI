#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CERT_DIR="$ROOT/secrets/certs"

mkdir -p "$CERT_DIR"

if [[ -f "$CERT_DIR/tls.crt" && -f "$CERT_DIR/tls.key" ]]; then
  echo "Local TLS certs already exist at $CERT_DIR"
  exit 0
fi

openssl req -x509 -newkey rsa:4096 \
  -keyout "$CERT_DIR/tls.key" \
  -out "$CERT_DIR/tls.crt" \
  -days 365 -nodes \
  -subj "/CN=api.ai-finops.local/O=AI FinOps Local Dev" \
  -addext "subjectAltName=DNS:api.ai-finops.local,DNS:app.ai-finops.local,DNS:localhost,IP:127.0.0.1"

chmod 600 "$CERT_DIR/tls.key"
echo "Generated local TLS certs in $CERT_DIR"