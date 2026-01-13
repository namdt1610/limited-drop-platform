# AI Coding Assistant Instructions for Node Ecommerce Platform

## Architecture Overview

- **Backend**: Go 1.25 + Fiber v3 + SQLite 3 (Split Architecture: Reader/Writer pools)
- **Database**: SQLite with WAL mode + optimized PRAGMA (\_synchronous=NORMAL, \_busy_timeout=5000, \_foreign_keys=on)
  - Reader Pool: 100 concurrent connections (parallel reads)
  - Writer Pool: 1 serialized connection (ordered writes)
  - SmartExecutor: Automatic routing (Query→Reader, Exec→Writer)
- **Frontend**: Alpine.js v3.13.4 (alpine/) — Vite v7.3.0, Bun v1.0.0 for dev; Tailwind CSS
- **Key Features**: SYMBICODE anti-counterfeit, limited drop with pessimistic locking, admin QR generation
- **Deployment**: Docker Compose for dev, systemd for production
- **Performance**: Optimized for 8,000-10,000 RPS with <100ms p99 latency

## Developer Workflows

- **Start Dev**: `docker-compose up --build` (SQLite local, backend/frontend local)
- **Build Backend**: `make build-prod` (optimized), `make build-linux-upx` (compressed for VPS)
- **Load Test**: `k6 run k6-runner.js -e K6_MODE=kpi-10k -e RATE=10000 -e DURATION=30s` (benchmark at 10k RPS)
- **Debug**: `go tool pprof` for CPU/memory profiling; SQLite WAL mode enables concurrent debugging
- **Deploy**: `make start` (Docker) or systemd services in `config/systemd/`
- **Database**: No schema migrations needed (auto-migrated), WAL file synced safely

## Code Patterns & Conventions

- **Go Modules**: `router.go` (handlers), `service.go` (business logic), `queries.go` (raw SQL for perf)
  - Example: `internal/modules/product/router.go` handles HTTP requests, calls `GetProducts` from `service.go`, which uses raw SQL in `queries.go` for performance.
  - Database: SmartExecutor automatically routes SELECT to Reader pool, INSERT/UPDATE/DELETE to Writer
- **Frontend SRP**: Feature-specific code in `src/features/`, shared in `src/shared/`
  - Example: `src/features/products/` contains components, hooks, services, screens for product-related UI.
- **Performance**: Use indexed queries to avoid N+1, proper indexes (B-tree/partial/composite); raw SQL for read-heavy ops
  - SQLite with WAL mode: Readers never blocked by writers, 100 parallel readers + 1 serialized writer
  - Read operations: ~50-200µs (microseconds), Write operations: ~100-400µs (serialized but fast)
- **Error Handling**: Fail fast, clear error messages, Vietnamese comments in code.
  - Example: Return `c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})` in handlers.
- **Security**: Rate limiting disabled in dev, JWT auth, input validation, foreign key constraints enforced

## Key Files & Directories

- `backend/main.go`: Entry point, middleware setup (compression, CORS, health check); auto-migrates models with GORM.
- `backend/internal/database/`: Split Architecture implementation (database.go, smart_executor.go)
  - `database.go`: DBInstance with Writer (1 conn) + Reader (100 conns)
  - `smart_executor.go`: Routes Query/QueryRow→Reader, Exec/Begin→Writer
- `backend/config/config.go`: Database configuration with MaxWriteConns, MaxReadConns, BusyTimeout
- `backend/internal/modules/`: Feature modules (product, auth, limiteddrop, symbicode); each with router.go, service.go, queries.go.
- `frontend/src/features/`: SRP feature code, e.g., `products/` with components/, hooks/, services/.
- `docs/OPTIMIZATION_LOG.md`: Before/after metrics (180x latency improvement)
- `docs/TECHNICAL_DETAILS.md`: Deep dive on Split Architecture + WAL mode
- `docs/QUICK_REFERENCE.md`: Cheat sheet for performance optimization

## Integration Points

- **Google Drive**: Uploads via `google.golang.org/api` (gdrive-service-account.json).
- **Health Checks**: `/health` endpoint with DB checks.
- **Compression**: Gzip/Brotli enabled, ETags for caching.
- **Database**: SmartExecutor auto-routes queries (no code changes needed)

Prioritize performance optimizations, data-driven decisions, and maintain brutal simplicity in all changes.

LAW: Seed fake drop mỗi đợt chỉ drop từ 5 đến 10 sản phẩm (Scarity), ko có wave gì hết, thả 1 lần duy nhất. (Each drop has 1 signature item (alpha items), others are rare (beta items)), thời gian diễn ra drop là 15p, xong thì phải check stock ko đc bán lố (overselled), check order đc tạo, nói chung là tuyệt đối toàn vẹn (Integrity), backend ko đc sập trong giờ peak, winners/losers will be recieved email has been setup in backend services. Everything is setup based on 7 deadly sins. Put the world in my cage. Do not use emoji in my sourcode!.
I dont want ppl fake my pd, i wanna create a serial number (like public key) and secret key will be stored in db and will be printed as QR code when packing the pd for delivery process in every pd. customer can verify my true pd by scanning the QR, redirect to /verify page. if the code has never been activated, show "XÁC NHẬN VẬT CHỦ. SYMBIONT ĐÃ ĐƯỢC KÍCH HOẠT VÀO LÚC [GIỜ/NGÀY]. CHÀO MỪNG ĐẾN VỚI Donald club." (Đồng thời update vào DB là mã này đã chết). Nếu mã (SYMBICODE) đã bị kích hoạt trước đó (Lũ fake copy lại mã)HIỆN THÔNG BÁO ĐỎ MÁU: "CẢNH BÁO. MÃ NÀY ĐÃ BỊ SỬ DỤNG. BẠN ĐANG CẦM TRÊN TAY ĐỐNG RÁC FAKE. VỨT NÓ ĐI." the code named SYMBICODE
scr code is still messy, pls check more carefully, categorize them, no need to write doc after every task.

TODO:
"Tạo thư mục tests & Viết Unit Test cho toàn bộ cái project Go Fiber này cho tao. Cover 100% case." (Nó viết lòi mắt, Ngài ngồi chơi).
"Comment giải thích từng dòng code như thể tao là đứa trẻ 5 tuổi." (Để sau này Ngài đọc lại còn hiểu).
"Viết Documentation API chuẩn OpenAPI/Swagger." (Cực lười làm tay, nhưng AI làm phút mốt).
"Dịch toàn bộ Error Message sang tiếng Việt phong cách Donald Vibe."

Checkout Flow:
Sử dụng sức mạnh của Golang Goroutine để xử lý bất đồng bộ (Async).
Bước 1: User bấm Checkout.
Bước 2: Go Fiber ghi đơn hàng vào SQLite (Bắt buộc thành công, Writer serialized).
Bước 3: Go Fiber trả về 200 OK cho User ngay lập tức (User sướng, app mượt).
Bước 4: Trước khi trả về, ông bắn ra một Goroutine (go func() { ... }()) để:
Đẩy data vào GSheet.
Gửi email
