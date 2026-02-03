.PHONY: run dev build build-prod test bench profile clean fmt lint deps update seed-battle test-battle audit-battle run-battle-test test-symbicode

# ==============================================================================
# DEVELOPMENT
# ==============================================================================

# Run server (default)
run:
	@echo "🚀 Starting server..."
	@go run cmd/server/main.go

# Development mode with auto-reload (requires air)
dev:
	@echo "🔄 Starting dev server (Air)..."
	@if [ -f ~/go/bin/air ]; then \
		~/go/bin/air; \
	else \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
		~/go/bin/air; \
	fi

# ==============================================================================
# BENCHMARK & PROFILING (Lab)
# ==============================================================================

# Run Database vs Redis Benchmark
bench:
	@echo "🧪 Running Benchmark Lab..."
	@./benchmarks/db_vs_redis/run.sh

# Run Load Test (K6)
k6:
	@echo "🔥 firing 2000 virtual users..."
	@k6 run ../tests/load/k6-loadtest.js

# Capture CPU Profile (30s)
profile:
	@echo "📊 Capturing CPU profile (30s)..."
	@echo "Open http://localhost:8081 when browser launches"
	@go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30

# ==============================================================================
# BUILD
# ==============================================================================

# Build dev binary
build:
	@echo "🔨 Building binary (bin/server)..."
	@go build -o bin/server cmd/server/main.go

# Build production binary (Optimized)
build-prod:
	@echo "📦 Building production binary (bin/server)..."
	@CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o bin/server cmd/server/main.go
	@ls -lh bin/server

# Build for Linux
build-linux:
	@echo "🐧 Building for Linux..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o bin/server-linux cmd/server/main.go

# ==============================================================================
# TESTING & QUALITY
# ==============================================================================

# Run all tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Header Tests (Legacy)
seed-battle:
	@go build -o seed-battle-drop ./cmd/seed-battle-drop && ./seed-battle-drop
test-battle:
	@cd scripts && ./test-payment-battle.sh
audit-battle:
	@go build -o audit-battle ./cmd/audit-battle && ./audit-battle
run-battle-test:
	@cd scripts && ./run-complete-battle-test.sh
test-symbicode:
	@cd scripts && ./test-symbicode.sh

# Format code
fmt:
	@go fmt ./...

# Lint code
lint:
	@golangci-lint run

# ==============================================================================
# UTILS
# ==============================================================================

clean:
	@rm -rf bin/ benchmark.db
	@go clean

deps:
	@go mod download
	@go mod tidy

update:
	@go list -u -m all
