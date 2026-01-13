# Backend - E-commerce API

High-performance e-commerce backend built with **Go**, **Fiber v3**, **GORM**, and **goccy/go-json**.

## Stack

| Component | Technology       | Purpose                        |
| --------- | ---------------- | ------------------------------ |
| Language  | Go 1.25+         | Performance, simplicity        |
| Framework | Fiber v3         | Express-like, fast             |
| SQL       | **GORM Raw SQL** | Raw SQL (performance-critical) |
| ORM       | GORM             | Admin CRUD only                |
| Database  | PostgreSQL 17    | Primary storage                |
| JSON      | goccy/go-json    | Zero-allocation JSON           |

**Frontend compatibility**:

- Alpine frontend (alpine/): Alpine.js v3.13.4, Vite v7.3.0, Bun v1.0.0

Keep frontend package.json in sync when bumping these versions.

## Quick Start

```bash
# Install dependencies
go mod download

# Create .env
cp .env.example .env
# Edit .env with your database credentials

# Run development server
make dev

# Or run directly
go run cmd/server/main.go
```

## Make Commands

```bash
make dev          # Development server with hot reload (air)
make build        # Dev build
make build-prod   # Production build (optimized, -33% size)
make build-linux  # Cross-compile for Linux VPS
make run          # Run binary
make test         # Run tests
make fmt          # Format code
make lint         # Lint code
make clean        # Clean build artifacts
make deps         # Install dependencies
```

## Architecture

### 3-Layer Architecture

```
internal/modules/{module}/
├── queries.go    # Data Access Layer (GORM queries)
├── service.go    # Business Logic Layer
└── router.go     # Presentation Layer (HTTP handlers)
```

### Project Structure

```
backend/
├── main.go                    # Entry point
├── Makefile                   # Build commands
├── internal/
│   ├── models/
│   │   └── models.go          # GORM models
│   ├── modules/
│   │   ├── product/           # Product module
│   │   ├── order/             # Order module
│   │   ├── cart/              # Cart module
│   │   ├── auth/              # Authentication
│   │   ├── user/              # User management (Admin)
│   │   ├── category/          # Categories
│   │   ├── review/            # Product reviews
│   │   ├── address/           # User addresses
│   │   ├── limiteddrop/       # Limited Drop endpoints
│   │   ├── analytics/         # Analytics endpoints
│   │   ├── media/             # File uploads
│   │   └── payment/           # Payment processing (PayOS)
│   ├── services/
│   │   ├── jwt.go             # JWT utilities
│   │   ├── jwt_blacklist.go   # JWT revocation
│   │   ├── auth_middleware.go # Auth + AdminOnly middleware
│   │   ├── payos.go           # PayOS integration
│   │   ├── cloudinary.go      # Image uploads
│   │   └── resend.go          # Email service
│   └── utils/
│       └── generics.go        # Go Generics utilities
├── cmd/
│   ├── seed/                  # Basic seeder (50 products)
│   └── seed-performance/      # Performance seeder (1K-50K)
└── docs/
    ├── postman-collection.json
    └── postman-environment.json
```

## API Endpoints

### Health

```
GET /health
```

### Authentication

```
POST /api/auth/register         # Register new user
POST /api/auth/login            # Login (returns JWT)
POST /api/auth/google           # Google OAuth
POST /api/auth/logout           # Logout (blacklist token)
POST /api/auth/refresh          # Refresh access token
GET  /api/auth/me               # Get current user
```

### Products (Public)

```
GET /api/products                      # List products (filters, search, pagination)
GET /api/products/:id                  # Product detail (active only)
GET /api/products/suggestions          # Search autocomplete
GET /api/products/variants/:id/stock   # Get variant stock
```

### Categories (Public)

```
GET /api/categories                    # List categories
GET /api/categories/:id                # Category detail
GET /api/categories/slug/:slug         # Category by slug
```

### Cart (Authenticated)

```
GET    /api/cart                       # Get cart
POST   /api/cart                       # Add to cart
PATCH  /api/cart/items/:id             # Update quantity
DELETE /api/cart/items/:id             # Remove item
DELETE /api/cart                       # Clear cart
```

### Orders

```
POST /api/orders/checkout              # Guest checkout (rate limited: 5/min)
POST /api/orders/checkout/auth         # Authenticated checkout
POST /api/orders/lookup                # Guest order lookup (order_number + email/phone)
```

### Reviews (Public Read)

```
GET  /api/reviews                      # List reviews
GET  /api/reviews/:id                  # Review detail
GET  /api/reviews/stats                # Review stats for product
POST /api/reviews                      # Create review (authenticated)
```

### Addresses (Authenticated)

```
GET    /api/addresses                  # List user addresses
POST   /api/addresses                  # Create address
PUT    /api/addresses/:id              # Update address
DELETE /api/addresses/:id              # Delete address
```

### Drops (Public)

```
GET  /api/drops                        # List active drops
GET  /api/drops/:id/status             # Drop status
POST /api/drops/:id/purchase           # Create payment link (authenticated); order is created on successful payment (first-to-pay wins)
```

### Orders (User Tracking)

```
GET  /api/orders                       # List orders by phone (query param: ?phone=123456789)
GET  /api/orders/:id                   # Get specific order by ID
```

### Analytics (Public Tracking)

```
POST /api/analytics/track                      # Track single event
POST /api/analytics/batch                      # Batch track events
```

### Payment

```
POST /api/payment/payos/checkout               # Create PayOS checkout
GET  /api/payment/payos/verify/:orderCode      # Verify payment status
POST /api/payment/payos/webhook                # PayOS webhook (updates order)
```

---

## Admin API Endpoints

All admin endpoints require authentication + admin role.

### Admin: Users

```
GET    /api/admin/users                        # List all users
GET    /api/admin/users/:id                    # User detail
POST   /api/admin/users                        # Create user
PATCH  /api/admin/users/:id                    # Update user
DELETE /api/admin/users/:id                    # Delete user
PATCH  /api/admin/users/:id/toggle-active      # Lock/Unlock user
```

### Admin: Products

```
GET    /api/admin/products                     # List ALL (incl. inactive)
POST   /api/admin/products                     # Create product
PUT    /api/admin/products/:id                 # Update product
DELETE /api/admin/products/:id                 # Delete product
PATCH  /api/admin/products/:id/toggle-active   # Toggle active
POST   /api/admin/products/:id/variants        # Create variant
PUT    /api/admin/products/variants/:id        # Update variant
DELETE /api/admin/products/variants/:id        # Delete variant
PATCH  /api/admin/products/variants/:id/stock  # Update stock
```

### Admin: Categories

```
POST   /api/admin/categories                   # Create category
PATCH  /api/admin/categories/:id               # Update category
DELETE /api/admin/categories/:id               # Delete category
```

### Admin: Orders

```
GET   /api/admin/orders                        # List all orders
GET   /api/admin/orders/:id                    # Order detail
PATCH /api/admin/orders/:id/status             # Update order status
```

### Admin: Reviews

```
PATCH  /api/admin/reviews/:id                  # Update review
DELETE /api/admin/reviews/:id                  # Delete review
```

```

```

### Admin: Drops

```
POST /api/admin/drops                  # Create drop
```

### Admin: Analytics

```
GET /api/admin/analytics/funnel                          # Funnel stats
GET /api/admin/analytics/trending                        # Trending products
GET /api/admin/analytics/segments/viewers-not-purchased  # Retargeting segment
```

---

## Security Features

### 1. Rate Limiting

```go
// Global: 60 requests/minute per IP
app.Use(services.RateLimitMiddleware(services.DefaultRateLimitConfig))

// Checkout: 5 requests/minute per IP (stricter)
app.Post("/api/orders/checkout",
    services.RateLimitMiddleware(CheckoutRateLimitConfig),
    publicCheckout,
)
```

### 2. Checkout Idempotency

Prevent double-click duplicate orders:

```bash
# Client sends idempotency key
curl -X POST /api/orders/checkout \
  -H "X-Idempotency-Key: uuid-here" \
  -d '...'

# If same key sent again, returns cached order (no duplicate)
```

### 3. AdminOnly Middleware

```go
// All admin routes require auth + admin role
admin := app.Group("/api/admin/users",
    services.AuthMiddleware(db),
    services.AdminOnly,
)
```

### 4. JWT Blacklist

---

## Performance Features

### 1. ETag Caching

```go
// Automatic 304 Not Modified for unchanged responses
app.Use(etag.New())
```

### 2. Compression

```go
// Gzip/Brotli compression (90% smaller responses)
app.Use(compress.New(compress.Config{
    Level: compress.LevelBestSpeed,
}))
```

### 3. goccy/go-json

Zero-allocation JSON encoder, 2-3x faster than `encoding/json`.

### 5. Database Indexes

All frequently queried fields are indexed (see `models/models.go`).

### 6. Connection Pooling

```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(50)
sqlDB.SetConnMaxLifetime(24 * time.Hour)
```

---

## Environment Variables

```bash
# Database (required)
DATABASE_URL=postgresql://user:pass@localhost:5432/dbname

# Server
PORT=3030
GOGC=200


# JWT (required)
JWT_SECRET=your-secret-key-min-32-chars

# PayOS (required for payments)
PAYOS_CLIENT_ID=...
PAYOS_API_KEY=...
PAYOS_CHECKSUM_KEY=...

# Optional
CLOUDINARY_URL=cloudinary://...
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
RESEND_API_KEY=...
```

---

## Database Seeding

```bash
# Basic seed (50 products)
go run cmd/seed/main.go

# Performance seed
PRODUCTS_COUNT=1000 FORCE_SEED=true go run cmd/seed-performance/main.go
```

---

## License

MIT
