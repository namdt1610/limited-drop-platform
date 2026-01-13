#!/bin/bash
# Start NocoDB (with SQLite) and Backend + Frontend Alpine (local dev)

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"

YELLOW='\033[1;33m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

BACKEND_PID=""
FRONTEND_PID=""

cleanup() {
    echo ""
    echo -e "${YELLOW}Stopping dev servers...${NC}"

    # Kill air process (backend hot reload)
    if [ -n "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
        # Also kill any air processes
        pkill -f "air" 2>/dev/null || true
    fi

    [ -n "$FRONTEND_PID" ] && kill $FRONTEND_PID 2>/dev/null || true

    echo -e "${GREEN}[OK] Dev servers stopped${NC}"
    echo -e "${YELLOW}Note: Docker services still running. Stop with: docker compose down${NC}"
    exit 0
}

trap cleanup INT TERM

# Detect docker compose command
if docker compose version > /dev/null 2>&1; then
    DOCKER_COMPOSE="docker compose"
elif command -v docker-compose > /dev/null 2>&1; then
    DOCKER_COMPOSE="docker-compose"
else
    echo -e "${RED}[X] Docker Compose is not installed${NC}"
    exit 1
fi

# Use dev compose file
COMPOSE_FILE="-f $PROJECT_ROOT/docker-compose.dev.yml"

# Check Docker permission
if ! docker ps > /dev/null 2>&1; then
    echo -e "${RED}[X] Docker permission denied${NC}"
    echo -e "${YELLOW}Fix: sudo usermod -aG docker $USER && logout/login${NC}"
    exit 1
fi

# 1. Start Docker services
echo -e "${YELLOW}[1/5] Starting NocoDB with SQLite...${NC}"
cd "$PROJECT_ROOT"
$DOCKER_COMPOSE $COMPOSE_FILE up -d nocodb

echo -e "${YELLOW}Waiting for NocoDB...${NC}"
for i in $(seq 1 30); do
    if curl -s http://localhost:8080 > /dev/null 2>&1; then
        echo -e "${GREEN}[OK] NocoDB ready${NC}"
        break
    fi
    if [ "$i" -eq 30 ]; then
        echo -e "${RED}[X] NocoDB not ready (timeout)${NC}"
        exit 1
    fi
    sleep 1
done

# 2. Seed Database
echo -e "${YELLOW}[2/5] Seeding database...${NC}"
cd "$PROJECT_ROOT/backend"
if [ -f .env ]; then
    set -o allexport
    source .env
    set +o allexport
fi
if go run cmd/seed/main.go 2>&1 | grep -q "already seeded"; then
    echo -e "${GREEN}[OK] Database already seeded${NC}"
else
    echo -e "${GREEN}[OK] Database seeded${NC}"
fi

# 3. Start Backend
echo -e "${YELLOW}[3/5] Starting backend (Go with hot reload)...${NC}"
cd "$PROJECT_ROOT/backend"
make dev-server &
BACKEND_PID=$!
echo -e "${GREEN}[OK] Backend started with hot reload (PID: $BACKEND_PID)${NC}"

# 4. Start Frontend Alpine
echo -e "${YELLOW}[4/5] Starting frontend-alpine (Alpine.js)...${NC}"
cd "$PROJECT_ROOT/alpine"

# Use Bun exclusively for frontend dev
BUN_CMD="bun"
if ! command -v bun > /dev/null 2>&1; then
    if [ -x "$HOME/.bun/bin/bun" ]; then
        BUN_CMD="$HOME/.bun/bin/bun"
    else
        echo -e "${RED}[X] Bun is required for frontend dev. Install Bun: https://bun.sh${NC}"
        kill $BACKEND_PID 2>/dev/null
        exit 1
    fi
fi

echo -e "${YELLOW}Installing alpine dependencies with Bun...${NC}"
$BUN_CMD install

echo -e "${YELLOW}Starting alpine with Bun + Vite + SWC...${NC}"
$BUN_CMD run dev &
FRONTEND_PID=$!
echo -e "${GREEN}[OK] Alpine frontend started (PID: $FRONTEND_PID)${NC}"

# 5. Done
echo ""
echo -e "${GREEN}[5/5] All services running!${NC}"
echo ""
echo "  Backend:       http://localhost:3030"
echo "  Frontend:      http://localhost:3000"
echo "  NocoDB:        http://localhost:8080"
echo "  Database:      SQLite (./database.db)"
echo ""
echo "Press Ctrl+C to stop dev servers"
echo ""

wait

