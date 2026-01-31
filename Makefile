# Makefile - Simplified
# Run `make` or `make help` to see available commands

.PHONY: help dev test build deploy clean

# =============================================================================
# HELP
# =============================================================================
help:
	@echo "Usage: make [command]"
	@echo ""
	@echo "Development:"
	@echo "  dev          Start dev environment (docker)"
	@echo "  dev-local    Start backend + frontend locally"
	@echo "  logs         View docker logs"
	@echo ""
	@echo "Testing:"
	@echo "  test         Run all tests"
	@echo "  test-be      Run backend tests only"
	@echo "  test-fe      Run frontend tests only"
	@echo "  test-load    Run K6 load test"
	@echo "  test-integrity Run Stress + Integrity Verifiction test"
	@echo ""
	@echo "Build:"
	@echo "  build        Build for production"
	@echo ""
	@echo "Deploy:"
	@echo "  deploy       Deploy to production"
	@echo ""
	@echo "Database:"
	@echo "  db           Open SQLite shell"
	@echo "  db-backup    Backup database"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean        Stop and cleanup docker"

# =============================================================================
# DEVELOPMENT
# =============================================================================
dev:
	docker compose -f deploy/docker/docker-compose.dev.yml up --build

dev-local:
	@echo "Starting backend..."
	cd backend && go run ./cmd/server &
	@echo "Starting frontend..."
	cd alpine && bun run dev

server:
	cd backend && go run ./cmd/server

seed:
	cd backend && FORCE_SEED=true go run ./cmd/seed

test-k6:
	k6 run tests/load/k6-loadtest.js

dev-localstack:
	docker compose -f deploy/docker/docker-compose.localstack.yml up

logs:
	docker compose logs -f

stop:
	docker compose down

# =============================================================================
# TESTING
# =============================================================================
test: test-be test-fe
	@echo "All tests passed!"

test-be:
	cd backend && go test ./tests/unit/... -v

test-fe:
	cd alpine && bun run test:run

test-load:
	k6 run tests/load/k6-runner.js

test-integrity:
	chmod +x tests/load/verify-integrity.sh
	./tests/load/verify-integrity.sh

test-localstack:
	./deploy/scripts/test-localstack.sh

# =============================================================================
# BUILD
# =============================================================================
build: build-be build-fe
	@echo "Build complete!"

build-be:
	cd backend && CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/server

build-fe:
	cd alpine && bun run build

# =============================================================================
# DEPLOY
# =============================================================================
deploy:
	./deploy/scripts/deploy.sh

# =============================================================================
# DATABASE
# =============================================================================
db:
	sqlite3 ./backend/database/database.db

db-backup:
	./deploy/scripts/backup.sh

# =============================================================================
# CLEANUP
# =============================================================================
clean:
	docker compose down -v
	rm -rf backend/server alpine/dist

clean-all:
	docker compose down -v --rmi all
	docker system prune -f
