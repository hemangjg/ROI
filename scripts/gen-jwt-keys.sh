#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
KEY_DIR="$ROOT/secrets/jwt"

mkdir -p "$KEY_DIR"

if [[ -f "$KEY_DIR/private.pem" && -f "$KEY_DIR/public.pem" ]]; then
  echo "JWT keys already exist at $KEY_DIR"
  exit 0
fi

openssl genrsa -out "$KEY_DIR/private.pem" 2048
openssl rsa -in "$KEY_DIR/private.pem" -pubout -out "$KEY_DIR/public.pem"
chmod 600 "$KEY_DIR/private.pem"

echo "Generated JWT keys in $KEY_DIR"