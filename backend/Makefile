.PHONY: build build-prod build-upx build-backup run dev dev-server test clean





























































































































































```rm /tmp/test.dbsqlite3 /tmp/test.db ".schema"gunzip backups/database_2026-01-05_02-00-00.db.gz -c > /tmp/test.db# Or restore to temporary locationgunzip -c backups/database_2026-01-05_02-00-00.db.gz | sqlite3 "SELECT COUNT(*) FROM products;"# Extract and check database```bashTest backup integrity:## Verification```// Use github.com/aws/aws-sdk-go-v2/service/s3// In cmd/backup/main.go, add S3 upload after successful backup```goExtend backup tool to upload to S3:### Remote Backups (S3)```sudo systemctl restart ecommerce-backup.timersudo systemctl daemon-reload```bashThen reload:```OnCalendar=*-*-* 03:30:00  # 3:30 AM instead[Timer]```iniEdit `ecommerce-backup.timer`:### Different Backup TimeThen modify cleanup in code (default: 30 days)```  -compress=true  -output /backups/ecommerce \  -db /opt/node-ecommerce/database.db \ExecStart=/opt/node-ecommerce/bin/backup \```iniEdit `ecommerce-backup.service`:### Custom Retention Period## Advanced Options- **Daily backups**: ~150-900MB disk space for 30 days- **Compressed**: 5-30MB (~85% reduction)- **Uncompressed**: 50-200MB (depends on data)SQLite with gzip compression typically achieves:## Backup Size```journalctl -u ecommerce-backup.service --since today# Check logswatch -n 5 'ls -lh /backups/ecommerce/ | head -10'# Monitor realtimels -lh /backups/ecommerce/```bashCheck backup creation:## Monitoring- **Format**: `database_YYYY-MM-DD_HH-MM-SS.db.gz`- **Location**: `/backups/ecommerce/`- **Retention**: Last 30 days automatically- **Frequency**: Daily at 2:00 AM (UTC)## Backup Schedule```sudo systemctl start ecommerce.service# Restart servercp backups/database_2026-01-05_02-00-00.db database.db# Or from uncompressed backupgunzip backups/database_2026-01-05_02-00-00.db.gz -c > database.db# Restore from compressed backupsudo systemctl stop ecommerce.service# Stop server```bash## Restore from Backup```gunzip backups/database_2026-01-05_02-00-00.db.gz# Unzip backup if needed./scripts/backup.shchmod +x scripts/backup.sh# Or use script./bin/backup -db database.db -output ./backups -compress=true# Build backup tool```bash## Manual Backup```journalctl -u ecommerce-backup.service -f# View logssudo systemctl list-timers ecommerce-backup.timersudo systemctl status ecommerce-backup.timer# Check statussudo systemctl start ecommerce-backup.timersudo systemctl enable ecommerce-backup.timersudo systemctl daemon-reload# Enable and startsudo cp config/systemd/ecommerce-backup.timer /etc/systemd/system/sudo cp config/systemd/ecommerce-backup.service /etc/systemd/system/# Copy files```bash### 3. Install Systemd Timer```chmod 755 /backups/ecommercemkdir -p /backups/ecommerce```bash### 2. Create Backup Directory```make build-backupcd backend```bash### 1. Build Backup Binary## Setup- **Manual backup** via CLI- **Systemd timer** (daily at 2 AM)- **Automatic cleanup** (keep 30 days)- **Gzip compression** (~70% size reduction)Automated SQLite database backups with:## Overview# Build binary (development)
build:
	@echo "Building Go binary (dev)..."
	@go build -o bin/server main.go
	@echo "Build complete: bin/server"

# Build backup utility
build-backup:
	@echo "Building backup utility..."
	@go build -o bin/backup ./cmd/backup
	@echo "Build complete: bin/backup"

# Development server with hot reload (requires air)
dev-server:
	@echo "Starting development server with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Or run: make run"; \
	fi

# Build binary (production) - Optimized with stripped symbols
# -ldflags "-s -w": Remove debug info & symbol table (40-50% smaller)
# -trimpath: Remove file system paths from binary
build-prod:
	@echo "Building Go binary (production, optimized)..."
	@CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o bin/server main.go
	@echo "Build complete: bin/server"
	@ls -lh bin/server

# Build for Linux (deploy to VPS)
build-linux:
	@echo "Building for Linux amd64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o bin/server-linux main.go
	@echo "Build complete: bin/server-linux"
	@ls -lh bin/server-linux

# Build with UPX compression (50-70% smaller, slower startup)
# Install UPX: sudo apt install upx / brew install upx
build-upx:
	@echo "Building Go binary with UPX compression..."
	@CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o bin/server main.go
	@if command -v upx > /dev/null; then \
		echo "Compressing with UPX..."; \
		upx --best --lzma bin/server; \
	else \
		echo "UPX not installed. Install with: sudo apt install upx"; \
	fi
	@ls -lh bin/server

# Build for Linux with UPX (deploy to low-resource VPS)
build-linux-upx:
	@echo "Building for Linux amd64 with UPX..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o bin/server-linux main.go
	@if command -v upx > /dev/null; then \
		echo "Compressing with UPX..."; \
		upx --best --lzma bin/server-linux; \
	else \
		echo "UPX not installed. Install with: sudo apt install upx"; \
	fi
	@ls -lh bin/server-linux

# Run server
run:
	@echo "Starting server..."
	@go run cmd/server/main.go

# Development mode with auto-reload (requires air)
dev-server:
	@if [ -f ~/go/bin/air ]; then \
		~/go/bin/air; \
	else \
		echo "Installing air for hot reload..."; \
		go install github.com/air-verse/air@latest; \
		~/go/bin/air; \
	fi

# Development mode with auto-reload (legacy alias)
dev: dev-server

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean
	@echo "Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# Lint code
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Check for updates
update:
	@echo "Checking for dependency updates..."
	@go list -u -m all

# Fierce Battle Test Suite
.PHONY: seed-battle test-battle audit-battle run-battle-test

# Seed battle drop for testing
seed-battle:
	@echo "Seeding fierce battle drop..."
	@go build -o seed-battle-drop ./cmd/seed-battle-drop
	@./seed-battle-drop

# Run battle test only (requires existing drop)
test-battle:
	@echo "Running fierce battle test..."
	@cd scripts && ./test-payment-battle.sh

# Audit drop integrity
audit-battle:
	@echo "Auditing drop integrity..."
	@go build -o audit-battle ./cmd/audit-battle
	@./audit-battle

# Complete battle test (seed + test + audit)
run-battle-test:
	@echo "Running complete fierce battle test suite..."
	@cd scripts && ./run-complete-battle-test.sh

# SYMBICODE Anti-Counterfeit System
.PHONY: test-symbicode

# Test SYMBICODE anti-counterfeit system
test-symbicode:
	@echo "Testing SYMBICODE anti-counterfeit system..."
	@cd scripts && ./test-symbicode.sh
