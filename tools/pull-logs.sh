#!/bin/bash
set -euo pipefail

# Pull Docker Logs for MenUnderFire
# Usage: ./tools/pull-logs.sh [--since <duration>]
#
# Pulls recent logs from menunderfire-app and menunderfire-postgres containers.
# Writes logs to tools/logs/ directory and prints the file paths.

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
LOG_DIR="$PROJECT_ROOT/tools/logs"
SINCE="35s"

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --since)
      SINCE="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--since <duration>]"
      exit 1
      ;;
  esac
done

mkdir -p "$LOG_DIR"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
APP_LOG="$LOG_DIR/app_${TIMESTAMP}.log"
PG_LOG="$LOG_DIR/postgres_${TIMESTAMP}.log"

# Pull app container logs
if docker ps --format '{{.Names}}' | grep -q '^menunderfire-app$'; then
  docker logs --since "$SINCE" menunderfire-app > "$APP_LOG" 2>&1 || true
else
  echo "[$(date -Iseconds)] Container menunderfire-app is not running" > "$APP_LOG"
fi

# Pull postgres container logs
if docker ps --format '{{.Names}}' | grep -q '^menunderfire-postgres$'; then
  docker logs --since "$SINCE" menunderfire-postgres > "$PG_LOG" 2>&1 || true
else
  echo "[$(date -Iseconds)] Container menunderfire-postgres is not running" > "$PG_LOG"
fi

# Output file paths for consumption by analyze-logs.sh
echo "APP_LOG=$APP_LOG"
echo "PG_LOG=$PG_LOG"
