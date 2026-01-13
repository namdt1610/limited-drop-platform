#!/bin/bash

# Database Migration Helper Script
# Usage: ./run-migration.sh [option]
# Options:
#   safe    - Run safe migration (preserves data)
#   clean   - Run clean migration (destructive)
#   backup  - Create database backup

set -e

# Database connection details (modify as needed)
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-donald_local}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Database Migration Helper${NC}"
echo "================================="

# Function to create backup
create_backup() {
    echo -e "${YELLOW}Creating database backup...${NC}"
    BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
    PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME > $BACKUP_FILE
    echo -e "${GREEN}Backup created: $BACKUP_FILE${NC}"
}

# Function to run migration
run_migration() {
    local script=$1
    local description=$2

    echo -e "${YELLOW}Running $description...${NC}"
    echo -e "${RED}⚠️  Make sure you have a backup before proceeding!${NC}"
    read -p "Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Migration cancelled."
        exit 1
    fi

    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $script

    echo -e "${GREEN}Migration completed successfully!${NC}"
    echo -e "${YELLOW}Please restart your backend application to apply GORM AutoMigrate.${NC}"
}

# Main logic
case "$1" in
    "safe")
        run_migration "update-schema-to-models.sql" "Safe Migration (preserves data)"
        ;;
    "clean")
        echo -e "${RED}⚠️  WARNING: Clean migration will DROP tables and may lose data!${NC}"
        run_migration "migrate-to-simplified-schema.sql" "Clean Migration (destructive)"
        ;;
    "backup")
        create_backup
        ;;
    *)
        echo "Usage: $0 {safe|clean|backup}"
        echo ""
        echo "Commands:"
        echo "  safe   - Run safe migration that preserves existing data"
        echo "  clean  - Run clean migration (WARNING: destructive!)"
        echo "  backup - Create database backup"
        echo ""
        echo "Examples:"
        echo "  ./run-migration.sh backup"
        echo "  ./run-migration.sh safe"
        exit 1
        ;;
esac
