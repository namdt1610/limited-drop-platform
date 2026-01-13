# Limited Drop E-commerce Platform

A full-stack e-commerce platform specializing in limited-time drops with first-to-pay-wins mechanics, built with modern technologies and production-ready architecture.

## Quick Start

### Alpine Frontend Development (Current Focus)

```bash
# Install Bun if not already installed
curl -fsSL https://bun.sh/install | bash

# Start Alpine development server
cd alpine && bun install && bun run dev
```

Visit: http://localhost:3000

### Development (Full Stack)

```bash
# Start all services for development/testing
docker-compose -f docker-compose.dev.yml up --build
```

Visit: http://localhost:3000

### Production Deployment

```bash
# 1. Upload docker-compose.yml and run DB + NocoDB
docker-compose up -d

# 2. Build and deploy backend
cd backend && go build -o server ./cmd/server
# Deploy server binary with your preferred method

# 3. Deploy frontend to Vercel/Netlify
```

## Architecture Overview

### Tech Stack

**Frontend**: Alpine.js v3.13.4 + Vite v7.3.0 + Tailwind CSS + Franken UI

**Backend**: Go 1.25+ + Fiber v3 + GORM (Postgres in production, SQLite for local tests)
**Database**: PostgreSQL 17 (production), SQLite for local/test runs
**Content Management**: NocoDB (No-code admin panel)
**Image Storage**: Cloudinary
**Payment**: PayOS (QR payments with webhooks)
**Email**: Resend

> Note: Exact versions are read from the package.json file in `alpine/` — keep it in sync when bumping dependencies.

### Project Structure

```
node-ecommerce/
├── alpine/              # Alpine.js frontend (current focus)
├── backend/             # Go backend with Fiber
│   ├── cmd/server/      # Main server entry point
│   ├── internal/
│   │   ├── handlers/    # HTTP handlers
│   │   ├── service/     # Business logic
│   │   ├── repository/  # Data access layer
│   │   ├── models/      # Data models
│   │   └── integrations/# External service integrations
│   └── database/        # Database schema and migrations

├── docker-compose.yml   # Production: NocoDB + SQLite
├── docker-compose.dev.yml # Development: NocoDB only
├── Makefile             # Build and deployment automation
└── env-example.txt      # Environment variables template
```

### Core Business Logic

- **Limited Drop System**: First-to-pay-wins mechanics with real-time stock tracking
- **Race Condition Protection**: Database-level pessimistic locking prevents overselling
- **Payment Integration**: PayOS QR payments with webhook processing
- **Order Management**: Automatic order creation from successful payments
- **User Analytics**: Auto-create users from phone numbers for tracking
- **Email Notifications**: Winner/loser notifications with order details
- **SYMBICODE**: Anti-counterfeit verification system with unique QR codes

### Core Flow (Checkout → Payment → Fulfillment)

1. Client initiates checkout (POST /api/orders/checkout) — server writes order record to DB immediately and returns 200 + payment link (idempotency key supported to avoid duplicate orders).
2. Client completes payment via PayOS; PayOS sends a webhook (POST /api/payment/payos/webhook).
3. Webhook handler verifies signature and idempotency, marks order as PAID, _increments stock atomically_ (or adjusts sold counters), and queues async post-processing (emails, Google Sheets update) in a separate goroutine.
4. On successful payment, a Symbicode may be generated and printed as QR code for product verification; verifying a Symbicode updates activation state (preventing reuse).

Notes:

- Asynchronous side-effects (email, GSheet) run in background goroutines to keep user-facing latency low.
- All critical DB updates are wrapped in transactions and tested (see `tests/unit/drop` and `tests/unit/symbicode`).

### Run & Test

- Run backend tests:

```bash
cd backend && go test ./... -v
```

- Run Alpine frontend (dev):

```bash
cd alpine && bun install && bun run dev
# Visit http://localhost:3000
```

If you prefer npm/yarn, `npm install && npm run dev` will also work.

## API Reference

### Base URL

```
Production: https://yourdomain.com/api
Development: http://localhost:3030/api
```

### Health Check

```http
GET /health
```

Response:

```json
{
  "status": "ok",
  "message": "Service is healthy"
}
```

### Products

#### List All Products

```http
GET /api/products
```

Response:

```json
[
  {
    "id": 1,
    "name": "Product Name",
    "description": "Product description",
    "price": 100000,
    "thumbnail": "image_url",
    "images": ["image1.jpg", "image2.jpg"],
    "tags": ["tag1", "tag2"],
    "stock": 50,
    "is_active": 1,
    "status": 0
  }
]
```

#### Get Single Product

```http
GET /api/products/{id}
```

### Limited Drops

#### Get Active Drops

```http
GET /api/drops
```

Response:

```json
[
  {
    "id": 1,
    "product_id": 2,
    "start_time": "2025-12-29T10:00:00Z",
    "end_time": "2025-12-29T10:15:00Z",
    "name": "Limited Drop #1",
    "total_stock": 10,
    "drop_size": 5,
    "sold": 0,
    "is_active": 1
  }
]
```

#### Purchase Drop

```http
POST /api/drops/{id}/purchase
```

Request Body:

```json
{
  "quantity": 1,
  "name": "John Doe",
  "phone": "0123456789",
  "email": "john@example.com",
  "address": "123 Main St",
  "province": "Ho Chi Minh City",
  "district": "District 1",
  "ward": "Ward 1"
}
```

Response (Success):

```json
{
  "message": "Purchase initiated",
  "payment_url": "https://payos.vn/payment/...",
  "order_code": 123456789
}
```

Response (Error):

```json
{
  "error": "Drop sold out"
}
```

### PayOS Webhook

```http
POST /api/limited-drops/webhook/payos
```

Headers:

```
x-payos-signature: {signature}
```

Webhook Payload:

```json
{
  "code": "00",
  "desc": "success",
  "data": {
    "orderCode": 123456789,
    "amount": 100000,
    "status": "PAID",
    "description": "Limited Drop Purchase",
    "metadata": {
      "drop_id": "1",
      "product_id": "2",
      "customer_phone": "0123456789",
      "customer_email": "john@example.com",
      "customer_name": "John Doe",
      "shipping_address": "123 Main St, District 1, Ho Chi Minh City",
      "quantity": "1"
    },
    "paymentMethod": "QR"
  }
}
```

### Orders

#### Get Order by ID

```http
GET /api/orders/{id}
```

Response:

```json
{
  "id": 1,
  "total_amount": 100000,
  "created_at": "2025-12-29T10:05:00Z",
  "customer_phone": "0123456789",
  "shipping_address": {
    "address": "123 Main St",
    "province": "Ho Chi Minh City",
    "district": "District 1",
    "ward": "Ward 1"
  },
  "items": [
    {
      "product_id": 2,
      "name": "Product Name",
      "price": 100000,
      "quantity": 1
    }
  ],
  "payment_method": 1,
  "status": 2,
  "payos_order_code": 123456789
}
```

#### Get Orders by Phone

```http
GET /api/orders?phone=0123456789
```

Response:

```json
{
  "orders": [...],
  "count": 5
}
```

### SYMBICODE Verification

#### Verify SYMBICODE

```http
POST /api/symbicode/verify
```

Request Body:

```json
{
  "code": "550e8400-e29b-41d4-a716-446655440000"
}
```

Response (First Activation):

```json
{
  "symbicode": {
    "id": 1,
    "order_id": 123,
    "product_id": 2,
    "created_at": "2025-12-29T10:05:00Z",
    "code": "550e8400-e29b-41d4-a716-446655440000",
    "is_activated": 1,
    "activated_at": "2025-12-29T12:00:00Z",
    "activated_ip": "192.168.1.100"
  },
  "is_first_activation": true
}
```

Response (Already Activated):

```json
{
  "error": "SYMBICODE already activated"
}
```

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT,
    phone TEXT UNIQUE,
    total_spent INTEGER DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    last_purchase_at DATETIME,
    is_active INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Products Table

```sql
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    thumbnail TEXT,
    images TEXT DEFAULT '[]',
    tags TEXT DEFAULT '[]',
    price INTEGER DEFAULT 0,
    stock INTEGER DEFAULT 0,
    is_active INTEGER DEFAULT 1,
    status INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);
```

### Orders Table

```sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    total_amount INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    customer_phone TEXT,
    shipping_address TEXT,
    items TEXT DEFAULT '[]',
    payment_method INTEGER DEFAULT 0,
    status INTEGER DEFAULT 0,
    payos_order_code INTEGER UNIQUE
);
```

### Limited Drops Table

```sql
CREATE TABLE limited_drops (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    name TEXT NOT NULL,
    total_stock INTEGER NOT NULL DEFAULT 0,
    drop_size INTEGER NOT NULL DEFAULT 1,
    sold INTEGER NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### SYMBICODE Table

```sql
CREATE TABLE symbicode (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER,
    product_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    activated_at DATETIME,
    code TEXT NOT NULL UNIQUE,
    secret_key TEXT NOT NULL,
    activated_ip TEXT,
    is_activated INTEGER NOT NULL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Development Setup

### Prerequisites

- **Go 1.21+**: Backend development
- **Bun**: Frontend development (faster than npm/yarn)
- **Docker & Docker Compose**: Database and services
- **SQLite3**: Database client (optional)

### Environment Variables

Create `.env` file in backend directory:

```bash
# Database
DB_PATH=./database.db

# PayOS Integration
PAYOS_CLIENT_ID=your_payos_client_id
PAYOS_API_KEY=your_payos_api_key
PAYOS_CHECKSUM_KEY=your_payos_checksum_key

# Email (Resend)
RESEND_API_KEY=your_resend_api_key

# Cloudinary (Image Storage)
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret

# NocoDB
NOCODB_PUBLIC_URL=http://localhost:8080
NOCODB_JWT_SECRET=dev_jwt_secret
NOCODB_ADMIN_EMAIL=admin@example.com
NOCODB_ADMIN_PASSWORD=password
```

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run with hot reload
go run github.com/cosmtrek/air

# Or run directly
go run ./cmd/server

# Build for production
go build -o server ./cmd/server
```

### Frontend Development (Alpine.js)

```bash
cd alpine

# Install dependencies
bun install

# Start development server
bun run dev

# Build for production
bun run build
```

### Database Management

```bash
# Initialize database
cd backend
sqlite3 database.db < database/init.sql

# View database
sqlite3 database.db
.schema
.tables

# Backup database
cp database.db backup_$(date +%Y%m%d).db
```

## Deployment

### Production Architecture

- **NocoDB**: Docker container (SQLite database file)
- **Backend**: Go binary with systemd
- **Frontend**: CDN (Vercel/Netlify)

### Docker Deployment

```bash
# Production (NocoDB only)
docker-compose up -d

# Development (full stack)
docker-compose -f docker-compose.dev.yml up --build
```

### Manual Backend Deployment

```bash
# Build optimized binary
cd backend
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Deploy binary to server
scp server user@your-server:/path/to/app/

# Run with systemd or process manager
./server
```

### Environment Setup

```bash
# Copy environment template
cp env-example.txt .env

# Edit with production values
nano .env
```

## Security Features

### Payment Security

- **HMAC Signature Verification**: PayOS webhook signatures
- **Idempotency**: Duplicate webhook prevention via order codes
- **Metadata Validation**: Secure payment metadata handling

### Anti-Fraud Protection

- **Rate Limiting**: Request throttling
- **Honeypots**: Bot detection traps
- **Time Validation**: Purchase timing checks
- **IP Tracking**: SYMBICODE activation logging

### Data Protection

- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries
- **XSS Protection**: Input sanitization
- **CSRF Protection**: Secure token handling

## Business Logic

### Limited Drop Flow

1. **Drop Creation**: Admin creates limited drop with stock and time window
2. **Purchase Initiation**: Customer submits shipping info
3. **Payment Creation**: PayOS QR code generated with metadata
4. **Payment Processing**: Customer scans QR and pays
5. **Webhook Processing**: PayOS sends payment confirmation
6. **Order Creation**: System creates order and updates stock
7. **Notification**: Email sent to winner/loser

### Race Condition Prevention

- **Database Locking**: Pessimistic locking prevents overselling
- **Atomic Operations**: Stock decrement and order creation in transaction
- **Idempotency Keys**: PayOS order codes prevent duplicate processing
- **Rollback Mechanism**: Failed operations restore stock

### SYMBICODE System

- **Unique Codes**: UUID v7 for each product sale
- **QR Generation**: Printable codes for packaging
- **Verification Portal**: Web interface for authenticity checking
- **Activation Tracking**: First-use detection and IP logging
- **Anti-Counterfeit**: Prevents fake product circulation

## Testing

### Unit Tests

```bash
cd backend
go test ./tests/unit/...
```

### Integration Tests

```bash
cd backend
go test ./tests/integration/...
```

### Load Testing

```bash
# Install K6
# Run load tests (requires K6 scripts)
k6 run backend/k6/load-test-k6.js
```

## Performance

### Database Optimization

- **Indexes**: Optimized for query patterns
- **Connection Pooling**: Efficient database connections
- **Query Optimization**: Raw SQL for performance-critical paths
- **Memory Usage**: SQLite file-based (no server overhead)

### Caching Strategy

- **Browser Caching**: Static assets with proper headers
- **CDN**: Frontend assets distributed globally
- **Database Indexing**: Fast queries for hot paths

### Scalability

- **Stateless Backend**: Horizontal scaling ready
- **File-based Database**: Simple deployment
- **Webhook Processing**: Asynchronous payment handling
- **Email Queue**: Background notification processing

## Maintenance

### Database Backup

```bash
# Manual backup
cp database.db backup_$(date +%Y%m%d_%H%M%S).db

# Automated backup (cron)
0 2 * * * cp /path/to/database.db /path/to/backups/database_$(date +\%Y\%m\%d).db
```

### Log Management

```bash
# View application logs
docker-compose logs -f backend

# System monitoring
htop
df -h
du -sh database.db
```

### Health Checks

```bash
# API health
curl http://localhost:3030/health

# Database connectivity
sqlite3 database.db "SELECT COUNT(*) FROM products;"

# Service status
docker-compose ps
```

## Contributing

### Code Standards

- **SOLID Principles**: Clean architecture
- **Error Handling**: Comprehensive error management
- **Documentation**: Inline code comments
- **Testing**: Unit and integration test coverage

### Development Workflow

1. **Fork** the repository
2. **Create** feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** changes (`git commit -m 'Add amazing feature'`)
4. **Push** to branch (`git push origin feature/amazing-feature`)
5. **Open** Pull Request

## License

This project is proprietary software. All rights reserved.

## Support

For technical support or questions:

- Create an issue in the repository
- Check existing documentation
- Review API examples

---

**Version**: 2.0.0
**Last Updated**: December 29, 2025
**Status**: Production Ready
