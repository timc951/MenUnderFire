#!/bin/bash
set -euo pipefail

# ============================================================
# Crontab Setup for PostgreSQL Backups
#
# Schedule:
#   - Saturday 2:00 AM: Full backup (pg_basebackup)
#   - Sun-Fri  2:00 AM: Incremental backup (WAL switch)
#
# Run this script once to install the cron job.
# ============================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKUP_SCRIPT="$SCRIPT_DIR/pg-backup.sh"
LOG_FILE="/var/log/pg-backup.log"

# Ensure the backup script is executable
chmod +x "$BACKUP_SCRIPT"

# Create the log file if it doesn't exist
touch "$LOG_FILE" 2>/dev/null || {
    echo "WARNING: Cannot create $LOG_FILE. You may need to run this with sudo."
    echo "Alternatively, edit the LOG_FILE path in pg-backup.sh to a writable location."
}

# Build the cron entry - runs daily at 2:00 AM
CRON_ENTRY="0 2 * * * $BACKUP_SCRIPT"

# Check if the cron job already exists
if crontab -l 2>/dev/null | grep -qF "$BACKUP_SCRIPT"; then
    echo "Cron job already exists. Updating..."
    # Remove the old entry and add the new one
    (crontab -l 2>/dev/null | grep -vF "$BACKUP_SCRIPT"; echo "$CRON_ENTRY") | crontab -
else
    echo "Installing new cron job..."
    (crontab -l 2>/dev/null; echo "$CRON_ENTRY") | crontab -
fi

echo ""
echo "Cron job installed successfully:"
echo "  Schedule: Daily at 2:00 AM"
echo "  Saturday: Full backup (pg_basebackup)"
echo "  Sun-Fri:  Incremental backup (WAL archive)"
echo "  Log file: $LOG_FILE"
echo ""
echo "Current crontab:"
crontab -l 2>/dev/null | grep -F "$BACKUP_SCRIPT"
