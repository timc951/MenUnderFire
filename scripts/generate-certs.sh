#!/bin/bash
# Generate self-signed SSL certificates for PostgreSQL
# Usage: ./scripts/generate-certs.sh

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CERT_DIR="$PROJECT_ROOT/certs"

mkdir -p "$CERT_DIR"

openssl req -new -x509 -days 365 -nodes -text \
  -out "$CERT_DIR/server.crt" \
  -keyout "$CERT_DIR/server.key" \
  -subj "/CN=postgres"

chmod 600 "$CERT_DIR/server.key"

echo "SSL certificates generated in $CERT_DIR"
