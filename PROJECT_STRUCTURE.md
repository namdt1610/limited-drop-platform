# Limited Drop E-commerce Platform

## Project Structure

```
node-ecommerce/
|
|-- alpine/                 # Frontend (Alpine.js + Vite)
|   |-- src/
|   |   |-- components/     # UI components
|   |   |-- lib/            # Utilities (api, dom-utils, addresses)
|   |   |-- pages/          # Page logic
|   |   `-- main.js         # Entry point
|   |-- public/             # Static assets
|   `-- package.json
|
|-- backend/                # Backend (Go + Fiber)
|   |-- cmd/                # Entry points (server, seed, backup)
|   |-- internal/
|   |   |-- handlers/       # HTTP handlers
|   |   |-- service/        # Business logic
|   |   |-- repository/     # Data access
|   |   |-- models/         # Data models
|   |   `-- integrations/   # External services
|   |-- database/           # SQL schemas
|   `-- tests/              # Unit & integration tests
|
|-- deploy/                 # Deployment configs
|   |-- docker/             # Docker Compose files
|   |-- nginx/              # Nginx configs
|   |-- localstack/         # AWS LocalStack init scripts
|   `-- scripts/            # Deploy & utility scripts
|
|-- docs/                   # Documentation
|   |-- LUAN_VAN_TOT_NGHIEP.md   # Thesis document
|   |-- TECHNICAL_DETAILS.md     # Technical specs
|   `-- ...
|
|-- tests/                  # Cross-project tests
|   `-- load/               # K6 load testing
|
|-- config/                 # System configs
|   `-- systemd/            # Systemd service files
|
|-- .github/
|   `-- workflows/          # CI/CD pipelines
|       |-- ci.yml          # Basic CI (tests)
|       |-- ci-localstack.yml  # CI with LocalStack
|       |-- deploy.yml      # Manual deployment
|       `-- release.yml     # Release on tag
|
|-- .env.example            # Environment template
|-- Makefile                # Build commands
|-- README.md               # This file
`-- docker-compose*.yml     # Symlinks to deploy/docker/
```

## Quick Start

### Development

```bash
# Frontend
cd alpine && bun install && bun run dev

# Backend
cd backend && go run ./cmd/server

# Full stack with LocalStack
docker compose -f deploy/docker/docker-compose.localstack.yml up
```

### Testing

```bash
# Backend tests
cd backend && go test ./tests/unit/... -v

# Frontend tests
cd alpine && bun run test:run

# LocalStack integration
./deploy/scripts/test-localstack.sh
```

### Deployment

```bash
# Build production
cd backend && go build -o server ./cmd/server
cd alpine && bun run build

# Deploy (via GitHub Actions or manual)
./deploy/scripts/deploy.sh
```

## Key Features

- **Limited Drop System**: First-to-pay-wins with race condition protection
- **Atomic Stock Updates**: Database-level consistency
- **Symbicode**: Anti-counterfeit QR verification
- **PayOS Integration**: QR payment with webhooks

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Alpine.js, Vite, Tailwind CSS |
| Backend | Go 1.21+, Fiber v3 |
| Database | SQLite (dev), PostgreSQL (prod) |
| Payment | PayOS |
| Email | Resend |
| Storage | Cloudinary |
| CI/CD | GitHub Actions |
| Testing | Go test, Vitest, K6 |
