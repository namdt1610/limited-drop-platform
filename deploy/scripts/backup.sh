#!/bin/bash

# Database Backup Script
# Usage: ./backup.sh [output_dir] [compress]

set -e

DB_PATH="${1:-.}/database.db"
OUTPUT_DIR="${2:-.}/backups"
COMPRESS="${3:-true}"

# Check if backup binary exists
BACKUP_BIN="./bin/backup"
if [ ! -f "$BACKUP_BIN" ]; then
    echo "Error: Backup binary not found at $BACKUP_BIN"
    echo "Build it with: cd backend && go build -o bin/backup ./cmd/backup"
    exit 1
fi

# Run backup
$BACKUP_BIN -db "$DB_PATH" -output "$OUTPUT_DIR" -compress=$COMPRESS

echo "Backup complete!"
