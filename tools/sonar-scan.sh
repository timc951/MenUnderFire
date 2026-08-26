#!/bin/bash
set -euo pipefail

# SonarQube Scanner for MenUnderFire
# Usage: ./tools/sonar-scan.sh [--token <token>] [--skip-tests] [--frontend-only] [--backend-only] [--debug]
#
# Requires: sonar-scanner CLI
#   Install: npm install -g sonarqube-scanner
#       or:  brew install sonar-scanner

SONAR_URL="http://172.16.255.10:9002"
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SKIP_TESTS=false
SONAR_TOKEN=""
SCAN_FRONTEND=true
SCAN_BACKEND=true
DEBUG_MODE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --token)
      SONAR_TOKEN="$2"
      shift 2
      ;;
    --skip-tests)
      SKIP_TESTS=true
      shift
      ;;
    --frontend-only)
      SCAN_BACKEND=false
      shift
      ;;
    --backend-only)
      SCAN_FRONTEND=false
      shift
      ;;
    --debug)
      DEBUG_MODE=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--token <token>] [--skip-tests] [--frontend-only] [--backend-only] [--debug]"
      exit 1
      ;;
  esac
done

# Check for token in env or argument
if [[ -z "$SONAR_TOKEN" ]]; then
  SONAR_TOKEN="${SONAR_TOKEN:-${SONAR_LOGIN:-}}"
fi

if [[ -z "$SONAR_TOKEN" ]]; then
  echo "Error: SonarQube token required."
  echo "  Pass via --token <token> or set SONAR_TOKEN env variable."
  echo ""
  echo "  Generate a token at: ${SONAR_URL}/account/security"
  exit 1
fi

# Verify sonar-scanner is installed
if ! command -v sonar-scanner &> /dev/null; then
  echo "Error: sonar-scanner not found."
  echo "  Install: npm install -g sonarqube-scanner"
  exit 1
fi

cd "$PROJECT_ROOT"

# Output directory for reports
OUTPUT_DIR="$PROJECT_ROOT/securityScan"
mkdir -p "$OUTPUT_DIR"

SCANNER_DEBUG_FLAG=""
if [[ "$DEBUG_MODE" == true ]]; then
  SCANNER_DEBUG_FLAG="-X"
fi

echo "=== MenUnderFire SonarQube Scan ==="
echo "SonarQube URL: $SONAR_URL"
echo "Results dir:   $OUTPUT_DIR"
if [[ "$DEBUG_MODE" == true ]]; then
  echo "Debug mode:    ON"
fi
echo ""

# --- Frontend ---
if [[ "$SCAN_FRONTEND" == true ]]; then
  echo "========================================="
  echo "  Scanning: Frontend (TypeScript/React)"
  echo "========================================="

  if [[ "$SKIP_TESTS" == false ]]; then
    echo "--- Generating frontend coverage ---"
    if [[ -d "frontend/node_modules" ]]; then
      if ! (cd frontend && npx vitest run --coverage); then
        echo ""
        echo "ERROR: Frontend tests failed. Fix the failing tests before running the scan."
        exit 1
      fi
      if [[ -f "frontend/coverage/lcov.info" ]]; then
        echo "✓ Frontend coverage report generated: frontend/coverage/lcov.info"
      else
        echo "Warning: Frontend coverage report not found at frontend/coverage/lcov.info"
      fi
    else
      echo "Skipping coverage (run npm install in frontend/ first)"
    fi
    echo ""
  fi

  echo "--- Running sonar-scanner for frontend ---"
  sonar-scanner $SCANNER_DEBUG_FLAG \
    -Dsonar.token="$SONAR_TOKEN" \
    -Dsonar.projectBaseDir="$PROJECT_ROOT/frontend" \
    -Dsonar.working.directory="$OUTPUT_DIR/.scannerwork-frontend" \
    -Dsonar.typescript.lcov.reportPaths=coverage/lcov.info \
    -Dsonar.javascript.lcov.reportPaths=coverage/lcov.info

  echo ""
  echo "Frontend results: ${SONAR_URL}/dashboard?id=MenUnderFire-frontend"
  echo ""
fi

# --- Backend ---
if [[ "$SCAN_BACKEND" == true ]]; then
  echo "========================================="
  echo "  Scanning: Backend (Go)"
  echo "========================================="

  if [[ "$SKIP_TESTS" == false ]]; then
    echo "--- Generating Go coverage ---"
    if command -v go &> /dev/null; then
      if ! (cd backend_go && go test ./... -coverprofile=coverage.out -coverpkg=./...); then
        echo ""
        echo "ERROR: Backend tests failed. Fix the failing tests before running the scan."
        exit 1
      fi
      if [[ -f "backend_go/coverage.out" ]]; then
        TOTAL_LINES=$(wc -l < backend_go/coverage.out)
        COVERED_LINES=$(grep -c " 1$" backend_go/coverage.out || true)
        echo "✓ Backend coverage report generated: backend_go/coverage.out"
        echo "  Lines in report: $TOTAL_LINES  |  Lines with coverage: $COVERED_LINES"
        echo "  Sample paths:"
        head -4 backend_go/coverage.out | tail -3 | sed 's/^/    /'
        if [[ "$COVERED_LINES" -eq 0 ]]; then
          echo ""
          echo "WARNING: Coverage report has 0 covered lines. Tests may not be exercising source code."
        fi
      else
        echo "Warning: Backend coverage report not found at backend_go/coverage.out"
      fi
    else
      echo "Skipping coverage (go not found)"
    fi
    echo ""
  fi

  echo "--- Running sonar-scanner for backend ---"
  sonar-scanner $SCANNER_DEBUG_FLAG \
    -Dsonar.token="$SONAR_TOKEN" \
    -Dsonar.projectBaseDir="$PROJECT_ROOT/backend_go" \
    -Dsonar.working.directory="$OUTPUT_DIR/.scannerwork-backend" \
    -Dsonar.go.coverage.reportPaths=coverage.out

  echo ""
  echo "Backend results: ${SONAR_URL}/dashboard?id=MenUnderFire-backend"
  echo ""
fi

echo "=== All scans complete ==="
