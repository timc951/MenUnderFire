#!/bin/bash
set -euo pipefail

# Trivy Security Scanner for MenUnderFire
# Usage: ./tools/trivy-scan.sh [--frontend-only] [--backend-only] [--docker-only] [--all]
#
# Requires: trivy CLI (https://aquasecurity.github.io/trivy/)
#   Install: brew install trivy
#       or:  sudo apt-get install trivy
#
# Sends results to Trivy server and saves local reports.

TRIVY_SERVER="http://172.16.255.10:8085"
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUTPUT_DIR="$PROJECT_ROOT/securityScan"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

SCAN_FRONTEND=false
SCAN_BACKEND=false
SCAN_DOCKER=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --frontend-only)
      SCAN_FRONTEND=true
      shift
      ;;
    --backend-only)
      SCAN_BACKEND=true
      shift
      ;;
    --docker-only)
      SCAN_DOCKER=true
      shift
      ;;
    --all)
      SCAN_FRONTEND=true
      SCAN_BACKEND=true
      SCAN_DOCKER=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--frontend-only] [--backend-only] [--docker-only] [--all]"
      exit 1
      ;;
  esac
done

# Default to --all if nothing specified
if [[ "$SCAN_FRONTEND" == false && "$SCAN_BACKEND" == false && "$SCAN_DOCKER" == false ]]; then
  SCAN_FRONTEND=true
  SCAN_BACKEND=true
  SCAN_DOCKER=true
fi

# Verify trivy is installed
if ! command -v trivy &> /dev/null; then
  echo "Error: trivy not found."
  echo "  Install: brew install trivy"
  echo "       or: sudo apt-get install trivy"
  exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

cd "$PROJECT_ROOT"

echo "=== MenUnderFire Trivy Security Scan ==="
echo "Trivy Server: $TRIVY_SERVER"
echo "Results dir:  $OUTPUT_DIR"
echo ""

# --- Frontend filesystem scan ---
if [[ "$SCAN_FRONTEND" == true ]]; then
  echo "========================================="
  echo "  Scanning: Frontend (filesystem)"
  echo "========================================="

  REPORT="$OUTPUT_DIR/frontend_${TIMESTAMP}.json"
  REPORT_TABLE="$OUTPUT_DIR/frontend_${TIMESTAMP}.txt"

  trivy filesystem \
    --server "$TRIVY_SERVER" \
    --scanners vuln,secret,misconfig \
    --format json \
    --output "$REPORT" \
    "$PROJECT_ROOT/frontend"

  trivy filesystem \
    --server "$TRIVY_SERVER" \
    --scanners vuln,secret,misconfig \
    --format table \
    --output "$REPORT_TABLE" \
    "$PROJECT_ROOT/frontend"

  echo "Results saved to:"
  echo "  JSON:  $REPORT"
  echo "  Table: $REPORT_TABLE"
  echo ""
fi

# --- Backend filesystem scan ---
if [[ "$SCAN_BACKEND" == true ]]; then
  echo "========================================="
  echo "  Scanning: Backend (filesystem)"
  echo "========================================="

  REPORT="$OUTPUT_DIR/backend_${TIMESTAMP}.json"
  REPORT_TABLE="$OUTPUT_DIR/backend_${TIMESTAMP}.txt"

  trivy filesystem \
    --server "$TRIVY_SERVER" \
    --scanners vuln,secret,misconfig \
    --format json \
    --output "$REPORT" \
    "$PROJECT_ROOT/backend_go"

  trivy filesystem \
    --server "$TRIVY_SERVER" \
    --scanners vuln,secret,misconfig \
    --format table \
    --output "$REPORT_TABLE" \
    "$PROJECT_ROOT/backend_go"

  echo "Results saved to:"
  echo "  JSON:  $REPORT"
  echo "  Table: $REPORT_TABLE"
  echo ""
fi

# --- Docker image scans ---
if [[ "$SCAN_DOCKER" == true ]]; then
  echo "========================================="
  echo "  Scanning: Docker Images"
  echo "========================================="

  # App image (combined frontend + backend)
  APP_IMAGE="menunderfire-app"
  if docker image inspect "$APP_IMAGE" &> /dev/null 2>&1; then
    echo "--- Scanning image: $APP_IMAGE ---"
    REPORT="$OUTPUT_DIR/docker_app_${TIMESTAMP}.json"
    REPORT_TABLE="$OUTPUT_DIR/docker_app_${TIMESTAMP}.txt"

    trivy image \
      --server "$TRIVY_SERVER" \
      --scanners vuln,secret \
      --format json \
      --output "$REPORT" \
      "$APP_IMAGE"

    trivy image \
      --server "$TRIVY_SERVER" \
      --scanners vuln,secret \
      --format table \
      --output "$REPORT_TABLE" \
      "$APP_IMAGE"

    echo "Results saved to:"
    echo "  JSON:  $REPORT"
    echo "  Table: $REPORT_TABLE"
    echo ""
  else
    echo "Skipping $APP_IMAGE (image not found, run docker compose build first)"
    echo ""
  fi

  # Postgres image
  PG_IMAGE="postgres:16-alpine"
  echo "--- Scanning image: $PG_IMAGE ---"
  REPORT="$OUTPUT_DIR/docker_postgres_${TIMESTAMP}.json"
  REPORT_TABLE="$OUTPUT_DIR/docker_postgres_${TIMESTAMP}.txt"

  trivy image \
    --server "$TRIVY_SERVER" \
    --scanners vuln,secret \
    --format json \
    --output "$REPORT" \
    "$PG_IMAGE"

  trivy image \
    --server "$TRIVY_SERVER" \
    --scanners vuln,secret \
    --format table \
    --output "$REPORT_TABLE" \
    "$PG_IMAGE"

  echo "Results saved to:"
  echo "  JSON:  $REPORT"
  echo "  Table: $REPORT_TABLE"
  echo ""
fi

echo "=== All scans complete ==="
echo "All reports saved in: $OUTPUT_DIR"
