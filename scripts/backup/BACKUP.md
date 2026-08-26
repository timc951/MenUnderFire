# PostgreSQL Backup & Restore

## Overview

Backups run daily at 2:00 AM via cron:

| Day | Type | Method | What it captures |
|-----|------|--------|-----------------|
| Saturday | Full | `pg_basebackup` (compressed tar + WAL stream) | Complete database snapshot |
| Sun-Fri | Incremental | WAL archiving (write-ahead log segments) | All changes since last full backup |

Backups are stored in the `pgbackups` Docker volume at `/backups` inside the container.

## Setup

### 1. Recreate the postgres container

The container needs WAL archiving enabled. If it was already running before these changes:

```bash
docker compose -f /path/to/docker-compose.yml down
docker compose -f /path/to/docker-compose.yml up -d
```

Data is preserved — only the container is recreated, not the `pgdata` volume.

### 2. Install the cron job

```bash
chmod +x scripts/backup/pg-backup.sh scripts/backup/setup-crontab.sh
./scripts/backup/setup-crontab.sh
```

This installs a single cron entry that runs `pg-backup.sh` daily at 2:00 AM. The script detects the day of the week and runs the appropriate backup type automatically.

### 3. Verify

```bash
# Run a manual backup
./scripts/backup/pg-backup.sh

# Check the log
cat /var/log/pg-backup.log

# List backups inside the container
docker exec dev-postgres ls -la /backups/full/
docker exec dev-postgres ls -la /backups/wal/
```

## Backup Details

### Full backup (Saturday)

- Uses `pg_basebackup` with `--wal-method=stream` to capture a consistent snapshot
- Output format: compressed tar (`-Ft -z`)
- Stored at: `/backups/full/YYYY-MM-DD/`
- Each full backup records the latest WAL file name in `latest_wal.txt` for cleanup

### Incremental backup (Sun-Fri)

- PostgreSQL continuously archives WAL (write-ahead log) segments to `/backups/wal/`
- The daily incremental job forces a WAL segment switch (`pg_switch_wal()`) to ensure recent changes are flushed to the archive
- WAL files contain every database change since the last full backup

### Retention

- Full backups older than **4 weeks** are automatically deleted
- WAL files no longer needed by any remaining full backup are cleaned up
- To change retention, edit `RETENTION_WEEKS` in `pg-backup.sh`

## Restore

### Full restore (latest state)

```bash
# 1. Stop the application
docker compose down

# 2. Remove the existing data volume
docker volume rm dev-resources_pgdata

# 3. Start a fresh postgres container
docker compose up -d postgres

# 4. Wait for it to be ready, then stop it
docker compose stop postgres

# 5. Clear the fresh data directory and restore the base backup
docker run --rm \
  -v dev-resources_pgdata:/var/lib/postgresql/data \
  -v dev-resources_pgbackups:/backups \
  postgres:16-alpine sh -c "
    rm -rf /var/lib/postgresql/data/*
    tar -xzf /backups/full/YYYY-MM-DD/base.tar.gz -C /var/lib/postgresql/data/
    tar -xzf /backups/full/YYYY-MM-DD/pg_wal.tar.gz -C /var/lib/postgresql/data/pg_wal/
    cp /backups/wal/* /var/lib/postgresql/data/pg_wal/ 2>/dev/null || true
    touch /var/lib/postgresql/data/recovery.signal
    chown -R 70:70 /var/lib/postgresql/data
  "

# 6. Start postgres — it will replay WAL files automatically
docker compose up -d
```

Replace `YYYY-MM-DD` with the date of the full backup you want to restore from.

### Point-in-time recovery (PITR)

To restore to a specific point in time rather than the latest state, add a `recovery_target_time` to `postgresql.conf` before starting:

```bash
docker run --rm \
  -v dev-resources_pgdata:/var/lib/postgresql/data \
  postgres:16-alpine sh -c "
    echo \"recovery_target_time = '2026-02-14 15:30:00'\" >> /var/lib/postgresql/data/postgresql.conf
    echo \"recovery_target_action = 'promote'\" >> /var/lib/postgresql/data/postgresql.conf
  "
```

Then start the container. PostgreSQL will replay WAL files only up to the specified time.

## Architecture

```
docker-compose.yml
  └── postgres container
        ├── /var/lib/postgresql/data  (pgdata volume - live database)
        └── /backups                  (pgbackups volume)
              ├── full/
              │     ├── 2026-02-08/   (Saturday full backup)
              │     │     ├── base.tar.gz
              │     │     ├── pg_wal.tar.gz
              │     │     └── latest_wal.txt
              │     └── 2026-02-15/   (next Saturday)
              └── wal/
                    ├── 000000010000000000000001
                    ├── 000000010000000000000002
                    └── ...           (continuous WAL archive)
```

## Configuration

| Setting | File | Default |
|---------|------|---------|
| Backup schedule | `setup-crontab.sh` | Daily at 2:00 AM |
| Retention period | `pg-backup.sh` (`RETENTION_WEEKS`) | 4 weeks |
| Container name | `pg-backup.sh` (`CONTAINER_NAME`) | `dev-postgres` |
| Log file | `pg-backup.sh` (`LOG_FILE`) | `/var/log/pg-backup.log` |
| WAL archive path | `postgres-ssl-setup.sh` (`archive_command`) | `/backups/wal/` |

## Troubleshooting

**Backup script fails with "container not found"**
- Verify the container is running: `docker ps | grep dev-postgres`
- Check the `CONTAINER_NAME` variable matches your `docker-compose.yml`

**No WAL files appearing in /backups/wal/**
- Verify WAL archiving is enabled: `docker exec dev-postgres psql -U postgres -c "SHOW archive_mode;"`
- Check archive command: `docker exec dev-postgres psql -U postgres -c "SHOW archive_command;"`
- Check for errors: `docker exec dev-postgres psql -U postgres -c "SELECT * FROM pg_stat_archiver;"`

**Full backup is very large**
- `pg_basebackup` captures the entire database cluster. This is expected.
- Adjust `RETENTION_WEEKS` if disk space is a concern.

**Cron job not running**
- Verify it's installed: `crontab -l`
- Check cron service is running: `systemctl status cron` or `service cron status`
- Check the log file: `cat /var/log/pg-backup.log`
