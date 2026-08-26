#!/bin/bash
set -euo pipefail

# Install/Uninstall MenUnderFire Log Monitor Cron Job
# Usage: ./tools/install-log-monitor.sh [--uninstall]
#
# Installs a crontab job that runs analyze-logs.sh every 30 seconds.
# Since cron's minimum interval is 1 minute, two entries are created:
#   - One at the top of each minute
#   - One with a 30-second sleep offset
#
# Prerequisites:
#   - docker CLI accessible from cron environment
#   - mail/sendmail/msmtp configured for email alerts (optional but recommended)

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SCRIPT_DIR="$PROJECT_ROOT/tools"
ANALYZE_SCRIPT="$SCRIPT_DIR/analyze-logs.sh"
LOG_DIR="$SCRIPT_DIR/logs"
MONITOR_LOG="$LOG_DIR/monitor.log"
CRON_MARKER="# menunderfire-log-monitor"

# Parse arguments
UNINSTALL=false
while [[ $# -gt 0 ]]; do
  case $1 in
    --uninstall)
      UNINSTALL=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--uninstall]"
      exit 1
      ;;
  esac
done

# --- Uninstall ---
if [[ "$UNINSTALL" == true ]]; then
  echo "Removing MenUnderFire log monitor cron jobs..."
  EXISTING=$(crontab -l 2>/dev/null || true)
  if echo "$EXISTING" | grep -q "$CRON_MARKER"; then
    echo "$EXISTING" | grep -v "$CRON_MARKER" | crontab -
    echo "Cron jobs removed."
  else
    echo "No MenUnderFire log monitor cron jobs found."
  fi
  exit 0
fi

# --- Install ---

# Check prerequisites
echo "=== MenUnderFire Log Monitor Installer ==="
echo ""

echo "Checking prerequisites..."

if ! command -v docker &> /dev/null; then
  echo "ERROR: docker not found. Install Docker first."
  exit 1
fi
echo "  [OK] docker"

if command -v mail &> /dev/null; then
  echo "  [OK] mail (mailutils)"
elif command -v sendmail &> /dev/null; then
  echo "  [OK] sendmail"
elif command -v msmtp &> /dev/null; then
  echo "  [OK] msmtp"
else
  echo "  [WARN] No mail command found (mail, sendmail, msmtp)."
  echo "         Email alerts will NOT be sent. Install one of these to enable alerts."
  echo "         e.g.: sudo apt-get install mailutils"
fi

# Verify scripts exist and are executable
if [[ ! -f "$ANALYZE_SCRIPT" ]]; then
  echo "ERROR: analyze-logs.sh not found at $ANALYZE_SCRIPT"
  exit 1
fi

chmod +x "$ANALYZE_SCRIPT"
chmod +x "$SCRIPT_DIR/pull-logs.sh"
echo "  [OK] Scripts are executable"

# Create log directory
mkdir -p "$LOG_DIR"

# Check if already installed
EXISTING=$(crontab -l 2>/dev/null || true)
if echo "$EXISTING" | grep -q "$CRON_MARKER"; then
  echo ""
  echo "MenUnderFire log monitor is already installed."
  echo "Use --uninstall to remove, then re-run to update."
  echo ""
  echo "Current cron entries:"
  echo "$EXISTING" | grep "$CRON_MARKER"
  exit 0
fi

# Build cron entries (every 30 seconds via two entries with sleep offset)
CRON_LINE_1="* * * * * ${ANALYZE_SCRIPT} >> ${MONITOR_LOG} 2>&1 ${CRON_MARKER}"
CRON_LINE_2="* * * * * sleep 30 && ${ANALYZE_SCRIPT} >> ${MONITOR_LOG} 2>&1 ${CRON_MARKER}"

# Append to existing crontab
{
  echo "$EXISTING"
  echo "$CRON_LINE_1"
  echo "$CRON_LINE_2"
} | crontab -

echo ""
echo "=== Installation Complete ==="
echo ""
echo "Cron jobs installed:"
echo "  $CRON_LINE_1"
echo "  $CRON_LINE_2"
echo ""
echo "Monitor log: $MONITOR_LOG"
echo "Incidents:   $LOG_DIR/incidents/"
echo "Alert email: timc950@hotmail.com"
echo ""
echo "To verify:   crontab -l | grep menunderfire"
echo "To remove:   $0 --uninstall"
