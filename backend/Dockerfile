# ============================================
# STAGE 1: Builder - Build the Go application
# ============================================
FROM golang:1.24-alpine AS builder

# Set GOTOOLCHAIN to use Go 1.25 if available
ENV GOTOOLCHAIN=go1.25.0

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make build-base

# Copy go mod file first
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy remaining files
COPY . .

# Tidy and verify dependencies
RUN go mod tidy && go mod verify

# Build the application with optimizations
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o ecommerce-backend \
    cmd/server/main.go

# ============================================
# STAGE 2: Runtime - Minimal production image
# ============================================
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies only
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/ecommerce-backend .

# Create non-root user for security
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser && \
    chown -R appuser:appuser /app

USER appuser

# Expose port
EXPOSE 3030

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3030/health || exit 1

# Run the application
CMD ["./ecommerce-backend"]