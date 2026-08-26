#!/bin/bash
set -uo pipefail

# MenUnderFire Log Anomaly Detection
# Usage: ./tools/analyze-logs.sh [--since <duration>] [--dry-run]
#
# Pulls recent Docker logs, scans for malicious patterns, and if threats are
# detected: sends an email alert and shuts down the app container.
#
# Designed to be run via cron every 30 seconds.
#
# Prerequisites:
#   - docker CLI accessible
#   - mail command (mailutils) or sendmail configured for email alerts

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SCRIPT_DIR="$PROJECT_ROOT/tools"
LOG_DIR="$SCRIPT_DIR/logs"
INCIDENT_DIR="$LOG_DIR/incidents"
COMPOSE_DIR="$PROJECT_ROOT"
ALERT_EMAIL="timc950@hotmail.com"
AUTH_FAILURE_THRESHOLD=10
SINCE="35s"
DRY_RUN=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --since)
      SINCE="$2"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--since <duration>] [--dry-run]"
      exit 1
      ;;
  esac
done

mkdir -p "$INCIDENT_DIR"

# --- Step 1: Pull logs ---
eval "$("$SCRIPT_DIR/pull-logs.sh" --since "$SINCE")"

# Verify log files exist
if [[ ! -f "$APP_LOG" || ! -f "$PG_LOG" ]]; then
  echo "[$(date -Iseconds)] ERROR: Failed to pull logs"
  exit 1
fi

# --- Step 2: Define detection patterns ---

# App container patterns (HTTP requests, access logs, application logs)
APP_PATTERNS=(
  # SQL Injection
  'UNION[[:space:]]+SELECT'
  'UNION[[:space:]]+ALL[[:space:]]+SELECT'
  "OR[[:space:]]+1[[:space:]]*=[[:space:]]*1"
  "OR[[:space:]]+'1'[[:space:]]*=[[:space:]]*'1'"
  'DROP[[:space:]]+TABLE'
  'DROP[[:space:]]+DATABASE'
  'INSERT[[:space:]]+INTO'
  'DELETE[[:space:]]+FROM'
  'UPDATE[[:space:]]+.*SET[[:space:]]'
  'SLEEP[[:space:]]*\('
  'BENCHMARK[[:space:]]*\('
  'INTO[[:space:]]+OUTFILE'
  'LOAD_FILE[[:space:]]*\('
  'INFORMATION_SCHEMA'
  'CHAR[[:space:]]*\([0-9]'

  # Path traversal
  '\.\.\/'
  '\.\.%2[fF]'
  '%2[eE]%2[eE]'
  '\.\.%5[cC]'
  '\/etc\/passwd'
  '\/etc\/shadow'
  '\/proc\/self'

  # XSS
  '<[[:space:]]*script'
  'javascript[[:space:]]*:'
  'onerror[[:space:]]*='
  'onload[[:space:]]*='
  '<[[:space:]]*iframe'
  '<[[:space:]]*svg[^>]*onload'
  'document\.cookie'
  'document\.location'

  # Command injection
  ';[[:space:]]*(ls|cat|id|whoami|wget|curl|nc|ncat|bash|sh|python|perl|ruby|php)[[:space:]]'
  '\|[[:space:]]*(cat|ls|id|bash|sh|whoami|wget|curl|nc)'
  '\$\(.*\)'
  '`[^`]*`'

  # Scanner / attack tool signatures
  'sqlmap'
  'nikto'
  'dirbuster'
  'gobuster'
  'wfuzz'
  'hydra'
  'masscan'
  'nessus'
  'burpsuite'
  'acunetix'
  'w3af'
  'zaproxy'
  'havij'

  # Shellshock
  '\(\)[[:space:]]*\{'

  # Log4Shell / JNDI injection
  '\$\{jndi:'
  '\$\{lower:'
  '\$\{upper:'

  # Server-Side Request Forgery indicators
  '169\.254\.169\.254'
  'metadata\.google\.internal'

  # Suspicious encoded payloads
  '%00'
  '%0[aAdD]'
)

# Postgres container patterns
PG_PATTERNS=(
  'FATAL.*password authentication failed'
  'FATAL.*no pg_hba.conf entry'
  'FATAL.*too many connections'
  'FATAL.*role.*does not exist'
  'FATAL.*database.*does not exist'
  'ERROR.*permission denied'
  'STATEMENT.*DROP'
  'STATEMENT.*TRUNCATE'
  'STATEMENT.*ALTER.*ROLE'
  'STATEMENT.*CREATE.*ROLE'
  'STATEMENT.*COPY.*TO'
  'STATEMENT.*pg_read_file'
  'STATEMENT.*pg_ls_dir'
)

# --- Step 3: Scan for anomalies ---

FINDINGS=""
THREAT_COUNT=0

# Build combined app pattern
APP_REGEX=""
for pattern in "${APP_PATTERNS[@]}"; do
  if [[ -n "$APP_REGEX" ]]; then
    APP_REGEX="$APP_REGEX|$pattern"
  else
    APP_REGEX="$pattern"
  fi
done

# Build combined postgres pattern
PG_REGEX=""
for pattern in "${PG_PATTERNS[@]}"; do
  if [[ -n "$PG_REGEX" ]]; then
    PG_REGEX="$PG_REGEX|$pattern"
  else
    PG_REGEX="$pattern"
  fi
done

# Scan app logs
APP_HITS=""
if [[ -s "$APP_LOG" ]]; then
  APP_HITS=$(grep -inE "$APP_REGEX" "$APP_LOG" 2>/dev/null || true)
fi

if [[ -n "$APP_HITS" ]]; then
  HIT_COUNT=$(echo "$APP_HITS" | wc -l)
  THREAT_COUNT=$((THREAT_COUNT + HIT_COUNT))
  FINDINGS="${FINDINGS}
=== APP CONTAINER THREATS (${HIT_COUNT} matches) ===
${APP_HITS}
"
fi

# Scan postgres logs
PG_HITS=""
if [[ -s "$PG_LOG" ]]; then
  PG_HITS=$(grep -inE "$PG_REGEX" "$PG_LOG" 2>/dev/null || true)
fi

if [[ -n "$PG_HITS" ]]; then
  HIT_COUNT=$(echo "$PG_HITS" | wc -l)
  THREAT_COUNT=$((THREAT_COUNT + HIT_COUNT))
  FINDINGS="${FINDINGS}
=== POSTGRES CONTAINER THREATS (${HIT_COUNT} matches) ===
${PG_HITS}
"
fi

# Check for excessive HTTP auth failures (brute force detection)
if [[ -s "$APP_LOG" ]]; then
  AUTH_FAILURES=$(grep -cE '"(401|403)"| 401 | 403 |status=401|status=403' "$APP_LOG" 2>/dev/null || echo "0")
  if [[ "$AUTH_FAILURES" -ge "$AUTH_FAILURE_THRESHOLD" ]]; then
    THREAT_COUNT=$((THREAT_COUNT + 1))
    FINDINGS="${FINDINGS}
=== BRUTE FORCE DETECTION ===
Excessive authentication failures detected: ${AUTH_FAILURES} in last ${SINCE}
Threshold: ${AUTH_FAILURE_THRESHOLD}
"
  fi
fi

# Check for excessive postgres auth failures
if [[ -s "$PG_LOG" ]]; then
  PG_AUTH_FAILURES=$(grep -c 'FATAL.*password authentication failed' "$PG_LOG" 2>/dev/null || echo "0")
  if [[ "$PG_AUTH_FAILURES" -ge 5 ]]; then
    THREAT_COUNT=$((THREAT_COUNT + 1))
    FINDINGS="${FINDINGS}
=== POSTGRES BRUTE FORCE DETECTION ===
Excessive postgres auth failures: ${PG_AUTH_FAILURES} in last ${SINCE}
"
  fi
fi

# --- Step 4: If threats found, alert and shutdown ---

if [[ "$THREAT_COUNT" -eq 0 ]]; then
  # Clean — remove temporary log files to avoid disk buildup
  rm -f "$APP_LOG" "$PG_LOG"
  exit 0
fi

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
INCIDENT_FILE="$INCIDENT_DIR/incident_${TIMESTAMP}.txt"

# Write incident report
cat > "$INCIDENT_FILE" <<EOF
===========================================
  MENUNDERFIRE SECURITY INCIDENT REPORT
===========================================
Timestamp:     $(date -Iseconds)
Threats Found: ${THREAT_COUNT}
Log Window:    Last ${SINCE}
App Log:       ${APP_LOG}
Postgres Log:  ${PG_LOG}

${FINDINGS}

=== RESPONSE ===
EOF

if [[ "$DRY_RUN" == true ]]; then
  echo "[DRY RUN] Would shut down menunderfire-app and send alert email"
  echo "$FINDINGS"
  echo "Action: DRY RUN — no shutdown or email sent" >> "$INCIDENT_FILE"
  exit 0
fi

# Send email alert
EMAIL_SUBJECT="[ALERT] MenUnderFire Security Incident - ${THREAT_COUNT} threats detected"
EMAIL_BODY="SECURITY ALERT — MenUnderFire

${THREAT_COUNT} potential malicious events detected at $(date -Iseconds).

The menunderfire-app container has been SHUT DOWN as a precaution.

${FINDINGS}

--- Incident Report ---
Saved to: ${INCIDENT_FILE}
App Log:  ${APP_LOG}
PG Log:   ${PG_LOG}

Review the incident report and logs, then restart the app when safe:
  cd ${COMPOSE_DIR} && docker compose up -d app
"

if command -v mail &> /dev/null; then
  echo "$EMAIL_BODY" | mail -s "$EMAIL_SUBJECT" "$ALERT_EMAIL"
  echo "Email sent via mail to ${ALERT_EMAIL}" >> "$INCIDENT_FILE"
elif command -v sendmail &> /dev/null; then
  {
    echo "Subject: $EMAIL_SUBJECT"
    echo "To: $ALERT_EMAIL"
    echo "Content-Type: text/plain; charset=UTF-8"
    echo ""
    echo "$EMAIL_BODY"
  } | sendmail "$ALERT_EMAIL"
  echo "Email sent via sendmail to ${ALERT_EMAIL}" >> "$INCIDENT_FILE"
elif command -v msmtp &> /dev/null; then
  {
    echo "Subject: $EMAIL_SUBJECT"
    echo "To: $ALERT_EMAIL"
    echo "Content-Type: text/plain; charset=UTF-8"
    echo ""
    echo "$EMAIL_BODY"
  } | msmtp "$ALERT_EMAIL"
  echo "Email sent via msmtp to ${ALERT_EMAIL}" >> "$INCIDENT_FILE"
else
  echo "WARNING: No mail command found (mail, sendmail, msmtp). Email NOT sent." >> "$INCIDENT_FILE"
  echo "[$(date -Iseconds)] WARNING: Could not send email alert — no mail utility found" >&2
fi

# Shut down the menunderfire app container (leave postgres for forensics)
echo "Shutting down menunderfire-app container..." >> "$INCIDENT_FILE"
cd "$COMPOSE_DIR" && docker compose stop app >> "$INCIDENT_FILE" 2>&1 || true
echo "App container stopped at $(date -Iseconds)" >> "$INCIDENT_FILE"

echo "[$(date -Iseconds)] ALERT: ${THREAT_COUNT} threats detected. App shut down. Incident: ${INCIDENT_FILE}"
