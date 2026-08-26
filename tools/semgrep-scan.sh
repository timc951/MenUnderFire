#!/bin/bash
set -euo pipefail

# Semgrep Scanner for MenUnderFire (runs via Docker)
# Usage: ./tools/semgrep-scan.sh [--frontend-only] [--backend-only] [--docker-only] [--all]
#
# Requires: docker

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUTPUT_DIR="$PROJECT_ROOT/securityScan"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
SEMGREP_IMAGE="semgrep/semgrep:latest"

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

# Verify docker is available
if ! command -v docker &> /dev/null; then
  echo "Error: docker not found."
  exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

cd "$PROJECT_ROOT"

echo "=== MenUnderFire Semgrep Scan ==="
echo "Image:       $SEMGREP_IMAGE"
echo "Results dir: $OUTPUT_DIR"
echo ""

# Helper: run semgrep via docker
# Usage: run_semgrep <mount_src> <container_path> <output_file> <extra_args...>
run_semgrep() {
  local mount_src="$1"
  local container_path="$2"
  local output_file="$3"
  local format="$4"
  shift 4

  docker run --rm \
    -v "${mount_src}:${container_path}:ro" \
    -v "${OUTPUT_DIR}:/output" \
    "$SEMGREP_IMAGE" \
    semgrep scan \
      --config auto \
      --${format} --output "/output/$(basename "$output_file")" \
      "$@" \
      "$container_path" || true
}

# --- Frontend ---
if [[ "$SCAN_FRONTEND" == true ]]; then
  echo "========================================="
  echo "  Scanning: Frontend (TypeScript/React)"
  echo "========================================="

  REPORT="$OUTPUT_DIR/semgrep_frontend_${TIMESTAMP}.json"
  REPORT_TEXT="$OUTPUT_DIR/semgrep_frontend_${TIMESTAMP}.txt"

  run_semgrep "$PROJECT_ROOT/frontend/src" "/src" "$REPORT" "json" \
    --include "*.ts" --include "*.tsx" --include "*.js" --include "*.jsx" \
    --exclude "node_modules" \
    --exclude "dist" \
    --exclude "*.test.*" \
    --exclude "*.spec.*" \
    --exclude "__tests__" \
    --exclude "__mocks__"

  run_semgrep "$PROJECT_ROOT/frontend/src" "/src" "$REPORT_TEXT" "text" \
    --include "*.ts" --include "*.tsx" --include "*.js" --include "*.jsx" \
    --exclude "node_modules" \
    --exclude "dist" \
    --exclude "*.test.*" \
    --exclude "*.spec.*" \
    --exclude "__tests__" \
    --exclude "__mocks__"

  echo "Results saved to:"
  echo "  JSON: $REPORT"
  echo "  Text: $REPORT_TEXT"
  echo ""
fi

# --- Backend ---
if [[ "$SCAN_BACKEND" == true ]]; then
  echo "========================================="
  echo "  Scanning: Backend (Go)"
  echo "========================================="

  REPORT="$OUTPUT_DIR/semgrep_backend_${TIMESTAMP}.json"
  REPORT_TEXT="$OUTPUT_DIR/semgrep_backend_${TIMESTAMP}.txt"

  run_semgrep "$PROJECT_ROOT/backend_go" "/src" "$REPORT" "json" \
    --include "*.go" \
    --exclude "vendor" \
    --exclude "*_test.go" \
    --exclude "*.exe" \
    --exclude "logs"

  run_semgrep "$PROJECT_ROOT/backend_go" "/src" "$REPORT_TEXT" "text" \
    --include "*.go" \
    --exclude "vendor" \
    --exclude "*_test.go" \
    --exclude "*.exe" \
    --exclude "logs"

  echo "Results saved to:"
  echo "  JSON: $REPORT"
  echo "  Text: $REPORT_TEXT"
  echo ""
fi

# --- Dockerfiles ---
if [[ "$SCAN_DOCKER" == true ]]; then
  echo "========================================="
  echo "  Scanning: Dockerfiles"
  echo "========================================="

  REPORT="$OUTPUT_DIR/semgrep_docker_${TIMESTAMP}.json"
  REPORT_TEXT="$OUTPUT_DIR/semgrep_docker_${TIMESTAMP}.txt"

  # Mount project root and scan only Dockerfiles
  docker run --rm \
    -v "${PROJECT_ROOT}:/src:ro" \
    -v "${OUTPUT_DIR}:/output" \
    "$SEMGREP_IMAGE" \
    semgrep scan \
      --config auto \
      --include "Dockerfile*" \
      --json --output "/output/$(basename "$REPORT")" \
      /src || true

  docker run --rm \
    -v "${PROJECT_ROOT}:/src:ro" \
    -v "${OUTPUT_DIR}:/output" \
    "$SEMGREP_IMAGE" \
    semgrep scan \
      --config auto \
      --include "Dockerfile*" \
      --text --output "/output/$(basename "$REPORT_TEXT")" \
      /src || true

  echo "Results saved to:"
  echo "  JSON: $REPORT"
  echo "  Text: $REPORT_TEXT"
  echo ""
fi

echo "=== All scans complete ==="
echo "All reports saved in: $OUTPUT_DIR"
