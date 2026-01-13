# Database Management Scripts

This directory contains scripts for managing the SQLite database used by the ecommerce platform.

## Scripts

### `db-backup.sh`

Creates a timestamped backup of the SQLite database.

```bash
# Manual backup
./scripts/db-backup.sh

# Or using Makefile
make db-backup
```

Backups are stored in the `./backups/` directory with format: `ecommerce_sqlite_YYYYMMDD_HHMMSS.db`

### `db-restore.sh`

Restores the database from a backup file.

```bash
# Restore from backup
./scripts/db-restore.sh backups/ecommerce_sqlite_20241211_120000.db

# Or using Makefile
make db-restore FILE=backups/ecommerce_sqlite_20241211_120000.db
```

**⚠️ WARNING**: This will overwrite the current database. A safety backup is created automatically before restoration.

### `db-shell`

Opens SQLite shell for direct database queries.

```bash
# Using Makefile
make db-shell
```

This opens the SQLite interactive shell where you can run SQL commands directly on the database.

## Database Location

- **Database file**: `./backend/database/database.db`
- **Backups directory**: `./backups/`

## Notes

- All scripts include error checking and user confirmation where appropriate
- Safety backups are created automatically during restore operations
- Scripts are designed to work with the Makefile targets for consistency
