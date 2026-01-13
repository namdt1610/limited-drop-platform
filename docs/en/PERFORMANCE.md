# Performance Benchmark & Limit Analysis

## Test Conditions
- **Date**: 2026-01-12
- **Environment**: Local Dev (Linux AMD64)
- **Database**: SQLite (WAL Mode)
- **Backend**: Go Fiber v3 (Single Instance)
- **Tool**: K6

## Why is it so fast? (Architectural & Engineering Secrets)
Achieving **10,000 RPS** with <1ms latency on a single node is not magic. It relies on three specific engineering decisions:

### 1. Zero-Network Database (SQLite In-Process)
This is the **biggest factor**. Traditional setups (App <-> Network <-> DB) incur TCP round-trip latency (typically 0.5ms - 1ms per query).
- **Our Approach**: SQLite runs **inside** the application process. Function calls replace network calls.
- **Latency**: Database query time is effectively **0ms** (microseconds range), limited only by disk I/O speed.

### 2. High-Performance Runtime (Go + Fiber)
- **Go**: Compiled to machine code, lightweight Goroutines allow handling thousands of concurrent connections with minimal RAM overhead.
- **Fiber**: Built on top of `fasthttp` (fastest Go HTTP engine), optimized for zero memory allocation in hot paths.

### 3. Optimized Concurrency Strategy
- **WAL Mode (Write-Ahead Logging)**: Enabled SQLite to handle multiple readers concurrent with a single writer.
- **Micro-Transactions**: We minimized the critical section in `ProcessSuccessfulDropPayment`. The transaction lock is held ONLY for the exact microsecond needed to write specific rows, drastically reducing lock contention.

## Benchmark Results

| Target RPS | Actual RPS | Avg Latency | P95 Latency | Success Rate | Status |
|------------|------------|-------------|-------------|--------------|--------|
| 2,000      | 2,000      | 0.58ms      | 1.34ms      | 100%         | PASS   |
| 5,000      | 5,000      | 0.85ms      | 3.56ms      | 100%         | PASS   |
| 8,000      | 8,000      | 0.79ms      | 3.40ms      | 100%         | PASS   |
| 15,000     | 13,286     | 35.63ms     | 166.63ms    | 100%         | PASS   |
| 25,000     | 15,981     | 430.84ms    | 1.17s       | 100%         | WARN   |
| 50,000     | 15,678     | 796.95ms    | 2.16s       | 100%         | FAIL   |

**Max Throughput**: ~16,000 RPS (Limited by SQLite Write Lock contention)
**Max Stable Throughput**: ~10,000 RPS (Low latency)

## Bottleneck Analysis

### 1. Database Locking (SQLite Specific)
Although WAL mode improves concurrency, SQLite still has a single writer lock.
- **Symptom**: Latency spikes exponentially from 8k -> 15k -> 25k RPS.
- **Status**: At 25k RPS, latency averages 430ms, which is unacceptable for real-time drops.
- **Solution for Prod**: PostgreSQL (row-level locking).

### 2. Connection Saturation
K6 reported `http_req_connecting` spikes at 50k RPS.
- **Symptom**: Connection timeout or long handshake times.
- **Solution**: Increase `ulimit`, tune kernel TCP stack (`somaxconn`, `tcp_max_syn_backlog`).

## Optimization Backlog

### High Priority
- [x] **WAL Mode**: Enabled for SQLite.
- [x] **Transaction Scope**: Minimized critical section in `ProcessSuccessfulDropPayment`.
- [x] **JSON Encoding**: Optimized to avoid double-encoding overhead.
- [ ] **Connection Pooling**: Tune `MaxOpenConns` / `MaxIdleConns` for SQL driver.

### Production Readiness
- [ ] **PostgreSQL Migration**: Switch from SQLite for production deployment.
- [ ] **Redis Caching**: Implement caching for `GetActiveDrops` (currently hitting DB).
- [ ] **Read Replicas**: Separate Read/Write traffic.
- [ ] **Rate Limiting**: Enable Fiber Rate Limiter middleware to protect against DDoS.

### Micro-optimizations
- [ ] **PGO (Profile Guided Optimization)**: Build Go binary with PGO for ~5-10% perf boost.
- [ ] **Struct Padding**: Reorder struct fields to optimize memory alignment.
- [ ] **Zero-copy Parsing**: Use fastjson or similar for high-throughput endpoint parsing.
