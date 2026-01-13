# E-commerce Platform

> **Version 2.0** - Production Ready

---

## **System Flow: Limited Drop Purchase**

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   FRONTEND      │ -> │  CHECKOUT MODAL  │ -> │  SHIPPING INFO  │
│   (Alpine)      │    │  (Form)          │    │  COLLECTION     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   BACKEND       │    │   PAYOS          │    │   METADATA       │
│   API CALL      │    │   PAYMENT LINK   │    │   STORAGE        │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   PAYOS QR      │    │   USER           │    │   PAYMENT        │
│   PAYMENT       │    │   SCAN & PAY     │    │   SUCCESS        │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   PAYOS         │    │   WEBHOOK        │    │   METADATA       │
│   WEBHOOK       │    │   CALLBACK       │    │   EXTRACTION     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   RACE          │    │   STOCK          │    │   USER           │
│   CONDITION     │    │   ATOMIC UPDATE  │    │   CREATION       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   ORDER         │    │   EMAIL          │    │   WINNER/        │
│   CREATION      │    │   NOTIFICATION   │    │   LOSER STATUS   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### **Key Data Flow**

1. **Frontend**: Collect shipping info in checkout modal
2. **Backend**: Store shipping info in PayOS payment metadata
3. **PayOS**: Process QR payment and return success webhook
4. **Webhook**: Extract shipping info from metadata
5. **Race Safe**: Atomic stock decrement with pessimistic locking
6. **Order**: Create order with shipping info, generate SYMBICODE
7. **Email**: Send winner notification or loser refund notification

## Architecture

```
node-ecommerce/
├── alpine/                # Alpine.js + Vite + Tailwind
│   └── src/features/      # Feature-based organization
├── backend/               # Go + Fiber v3 + GORM (SOLID)
│   └── internal/modules/  # Clean architecture modules
├── docs/                  # Complete documentation
├── scripts/               # Build & deploy scripts
└── docker-compose.yml     # PostgreSQL only
```

---

## Tech Stack

### Frontend

| Tech             | Purpose           |
| ---------------- | ----------------- |
| **Alpine.js 3**  | UI Framework      |
| **Vite 7**       | Build tool, HMR   |
| **Tailwind CSS** | Styling           |
| **Zod**          | Schema validation |

### Backend

| Tech              | Purpose                         |
| ----------------- | ------------------------------- |
| **Go 1.25**       | Language                        |
| **Fiber v3**      | Web framework (10x faster Node) |
| **GORM**          | ORM + Raw SQL                   |
| **PostgreSQL 17** | Primary database                |
| **PayOS SDK**     | Payment gateway integration     |
| **Resend**        | Email service                   |

---

## Features

### Core E-commerce

- **Product Catalog** - Static product data with bitwise flags
- **Limited Drop System** - First-to-pay-wins mechanics with race condition protection
- **Guest Checkout** - No account required for purchases
- **Order Tracking** - Base32 encoded order numbers (DV-{code}) + phone verification
- **User Creation** - Auto-create users from phone numbers for analytics
- **Email Notifications** - Winner/loser notifications via Resend

### Payment & Security

- **PayOS Integration** - QR payments with webhook processing
- **Metadata Storage** - Shipping info stored in PayOS payment metadata
- **Webhook Security** - HMAC signature verification (production)
- **Race Condition Safe** - Database pessimistic locking prevents overselling
- **Anti-bot Protection** - Honeypots, time traps, and rate limiting
- **Real-time Updates** - 2-second polling for live stock status

### Performance & Security

- **< 1.5s LCP** - Code splitting, ETag, compression
- **Rate Limiting** - 60/min global, 5/min checkout
- **Race Condition Safe** - Database pessimistic locking
- **~30KB bundle** - Alpine.js lightweight
- **SOLID Architecture** - Clean service separation

---

## Quick Start

### Requirements

- Go 1.25+
- Node.js 20+ or Bun
- PostgreSQL 17

### 1. Clone & Setup

```bash
git clone https://github.com/namdt1610/node-ecommerce.git
cd node-ecommerce
```

### 2. Start Database (Docker)

```bash
docker-compose up -d
```

### 3. Backend

```bash
cd backend
cp .env.example .env
# Edit .env with your database credentials

go mod download
make dev
```

### 4. Frontend

```bash
cd frontend
cp .env.example .env

bun install
bun run dev
```

### 5. Access

| Service      | URL                          |
| ------------ | ---------------------------- |
| Frontend     | http://localhost:3000        |
| Backend API  | http://localhost:3030        |
| Health Check | http://localhost:3030/health |
| PostgreSQL   | localhost:5432               |

---

## Project Structure

### Frontend (`/alpine`)

```
src/
├── components/         # Reusable UI components
│   ├── drop/           # Drop page components
│   └── landing/        # Landing page components
├── pages/              # Page logic files
│   ├── drop.js         # Drop page
│   ├── landing.js      # Landing page
│   └── ...
├── lib/                # API client, utilities
└── main.js             # Entry point
```

### Backend (`/backend`)

```
internal/
├── models/             # GORM models
├── modules/            # SOLID architecture (services + types)
│   ├── limiteddrop/    # First-to-pay-wins system
│   │   ├── router/     # HTTP routes, validation, handlers
│   │   ├── services/   # Business logic (SOLID services)
│   │   ├── types/      # Interfaces & DTOs
│   │   └── shared/     # Shared types & errors
│   ├── payment/        # PayOS webhook processing
│   │   ├── services/   # Payment & webhook services
│   │   ├── types/      # Payment interfaces
│   │   └── shared/     # Error handling
│   ├── product/        # Product management
│   └── ...
├── services/           # Shared services (PayOS, Resend, etc.)
└── shared/             # Common utilities
```

---

## Commands

### Backend

```bash
make dev          # Dev server (hot reload)
make build        # Dev build
make build-prod   # Production build (-33% size)
make test         # Run tests
```

### Frontend

```bash
bun dev           # Dev server
bun build         # Production build
bun preview       # Preview build
```

### Docker

```bash
docker-compose up -d      # Start services
docker-compose logs -f    # View logs
docker-compose down       # Stop services
```

---

## Documentation

| Doc                                       | Description               |
| ----------------------------------------- | ------------------------- |
| [01-QUICKSTART](./01-QUICKSTART.md)       | Getting started           |
| [02-ARCHITECTURE](./02-ARCHITECTURE.md)   | System design             |
| [03-DATABASE](./03-DATABASE.md)           | Database schema           |
| [04-PERFORMANCE](./04-PERFORMANCE.md)     | Optimization guide        |
| [05-FLASH-SALE](./05-FLASH-SALE.md)       | Flash sale implementation |
| [06-API](./06-API.md)                     | API reference             |
| [08-DOCKER-DEPLOY](./08-DOCKER-DEPLOY.md) | Deployment guide          |
| [PRD-CORE](./PRD-CORE.md)                 | Product requirements      |

---

## API Overview

### Public APIs

```
/api/products/*          # Product catalog
/api/limited-drops/*     # Limited drop operations
  ├── GET  /api/limited-drops/{id}     # Get drop details
  └── POST /api/limited-drops/{id}/purchase  # Purchase drop
/api/payment/track-order # Order lookup by Base32 order number + phone
/api/payment/payos/*     # PayOS operations
  ├── POST /api/payment/payos/checkout     # Create payment link
  ├── GET  /api/payment/payos/verify/{orderCode}  # Verify payment
  ├── POST /api/payment/payos/webhook      # Webhook receiver
```

### Protected (Auth Required)

```
/api/auth/*              # Authentication
/api/cart/*              # Cart
/api/addresses/*         # Addresses
```

```

---

## Performance Metrics

| Metric | Target | Achieved |
|--------|--------|----------|
| LCP | < 1.5s | ~1.2s |
| FCP | < 0.8s | ~0.6s |
| Bundle | < 150KB | ~120KB |
| Checkout API | < 700ms | ~400ms |

---

## Security Features

- Rate limiting (disabled in dev)
- Checkout idempotency (X-Idempotency-Key)
- JWT + Refresh tokens
- JWT Blacklist (logout revocation)
- CORS configured
- SQL injection prevention (GORM)

---

## License

MIT
```
