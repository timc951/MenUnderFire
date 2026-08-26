#!/bin/bash
set -euo pipefail

# ============================================================
# PostgreSQL Backup Script
#   - Saturday: Full backup via pg_basebackup
#   - Other days: Incremental backup via WAL archiving
#
# Requires WAL archiving enabled in PostgreSQL (see docker-compose.yml)
# ============================================================

CONTAINER_NAME="dev-postgres"
BACKUP_VOLUME_PATH="/backups"
RETENTION_WEEKS=4
DATE=$(date +%Y-%m-%d)
DAY_OF_WEEK=$(date +%u)  # 1=Monday, 6=Saturday, 7=Sunday
LOG_FILE="/var/log/pg-backup.log"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

check_container() {
    if ! docker inspect "$CONTAINER_NAME" &>/dev/null; then
        log "ERROR: Container '$CONTAINER_NAME' not found"
        exit 1
    fi
    if [ "$(docker inspect -f '{{.State.Running}}' "$CONTAINER_NAME")" != "true" ]; then
        log "ERROR: Container '$CONTAINER_NAME' is not running"
        exit 1
    fi
}

full_backup() {
    local backup_dir="$BACKUP_VOLUME_PATH/full/$DATE"

    log "Starting FULL backup to $backup_dir"

    # Create backup directory inside the container
    docker exec "$CONTAINER_NAME" mkdir -p "$backup_dir"

    # Run pg_basebackup (compressed tar format with checksum verification)
    docker exec "$CONTAINER_NAME" pg_basebackup \
        -U postgres \
        -D "$backup_dir" \
        -Ft \
        -z \
        -P \
        --checkpoint=fast \
        --wal-method=stream

    # Record which WAL file the backup ended at (for incremental restore)
    docker exec "$CONTAINER_NAME" sh -c \
        "psql -U postgres -t -c \"SELECT pg_walfile_name(pg_current_wal_lsn());\" > $backup_dir/latest_wal.txt"

    # Clean up WAL files that are older than the new full backup
    clean_old_wal "$backup_dir"

    log "FULL backup completed: $backup_dir"
}

incremental_backup() {
    log "Starting INCREMENTAL backup (WAL archive)"

    # Force PostgreSQL to switch to a new WAL segment so current changes are archived
    docker exec "$CONTAINER_NAME" psql -U postgres -c "SELECT pg_switch_wal();" >/dev/null

    # Verify WAL files are being archived
    local wal_count
    wal_count=$(docker exec "$CONTAINER_NAME" sh -c "ls $BACKUP_VOLUME_PATH/wal/ 2>/dev/null | wc -l")

    log "INCREMENTAL backup completed: $wal_count WAL files archived"
}

clean_old_wal() {
    local current_backup_dir="$1"

    # Find the oldest full backup we want to keep
    local cutoff_date
    cutoff_date=$(date -d "-${RETENTION_WEEKS} weeks" +%Y-%m-%d 2>/dev/null || date -v-${RETENTION_WEEKS}w +%Y-%m-%d)

    log "Cleaning backups older than $cutoff_date (retention: ${RETENTION_WEEKS} weeks)"

    # Remove old full backups
    docker exec "$CONTAINER_NAME" sh -c "
        for dir in $BACKUP_VOLUME_PATH/full/*/; do
            dirname=\$(basename \"\$dir\")
            if [ \"\$dirname\" '<' '$cutoff_date' ] 2>/dev/null; then
                rm -rf \"\$dir\"
                echo \"Removed old full backup: \$dirname\"
            fi
        done
    " 2>/dev/null | while read -r line; do log "$line"; done

    # Get the earliest WAL file needed by the oldest remaining full backup
    local earliest_wal
    earliest_wal=$(docker exec "$CONTAINER_NAME" sh -c "
        oldest_backup=\$(ls -d $BACKUP_VOLUME_PATH/full/*/ 2>/dev/null | head -1)
        if [ -n \"\$oldest_backup\" ] && [ -f \"\${oldest_backup}latest_wal.txt\" ]; then
            cat \"\${oldest_backup}latest_wal.txt\" | tr -d ' '
        fi
    ")

    # Remove WAL files older than what the oldest full backup needs
    if [ -n "$earliest_wal" ]; then
        docker exec "$CONTAINER_NAME" sh -c "
            for wal_file in $BACKUP_VOLUME_PATH/wal/*; do
                fname=\$(basename \"\$wal_file\")
                if [ \"\$fname\" '<' '$earliest_wal' ] 2>/dev/null; then
                    rm -f \"\$wal_file\"
                fi
            done
        "
        log "Cleaned WAL files older than $earliest_wal"
    fi
}

# ---- Main ----
log "=== Backup started ==="
check_container

if [ "$DAY_OF_WEEK" -eq 6 ]; then
    full_backup
else
    incremental_backup
fi

log "=== Backup finished ==="
