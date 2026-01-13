package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DBInstance chứa 2 con trỏ: 1 để Ghi, 1 để Đọc
// Chiến lược "Split Architecture" tối ưu hóa hiệu năng SQLite
type DBInstance struct {
	Writer *sql.DB // 1 connection duy nhất để tránh xung đột khóa
	Reader *sql.DB // Nhiều connections để đọc song song
}

var DB DBInstance

// Connect khởi tạo database với cấu hình Production Sweet Spot
func Connect(dbPath string) error {
	// DSN chuẩn cho Production
	// _journal_mode=WAL: Cho phép Đọc/Ghi song song
	// _synchronous=NORMAL: An toàn + nhanh (chỉ mất uncommitted transaction nếu mất điện)
	// _busy_timeout=5000: Chờ 5s trước khi báo lỗi (giảm error rate khi tải cao)
	// _foreign_keys=on: Bật foreign key constraints
	dsn := dbPath + "?_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=5000&_foreign_keys=on&cache=shared"

	var err error

	// 1. KHỞI TẠO WRITER (QUAN TRỌNG NHẤT)
	// Writer chỉ được phép có 1 Connection duy nhất để tránh xung đột khóa
	// Ép toàn bộ lệnh Ghi phải xếp hàng (Serialize) trong Go
	// Điều này nhanh hơn để SQLite tự lock file
	DB.Writer, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	// Cấu hình Writer: "Cổ chai" chủ động
	DB.Writer.SetMaxOpenConns(1) // Chỉ 1 connection duy nhất
	DB.Writer.SetMaxIdleConns(1) // Giữ 1 connection idle
	DB.Writer.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := DB.Writer.Ping(); err != nil {
		return err
	}

	// 2. KHỞI TẠO READER
	// Reader dùng chung file nhưng object *sql.DB khác
	// WAL mode cho phép đọc không chặn ghi
	DB.Reader, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	// Cấu hình Reader: Mở rộng theo CPU
	// Ví dụ: 100 kết nối đọc đồng thời
	DB.Reader.SetMaxOpenConns(100) // Hỗ trợ 100 concurrent reads
	DB.Reader.SetMaxIdleConns(100) // Keep 100 idle connections
	DB.Reader.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := DB.Reader.Ping(); err != nil {
		return err
	}

	log.Println("🚀 Database Connected: Production Sweet Spot Mode (WAL + Split Architecture)")
	return nil
}

// Close đóng cả Writer và Reader
func Close() error {
	if err := DB.Writer.Close(); err != nil {
		log.Printf("error closing Writer: %v", err)
	}
	if err := DB.Reader.Close(); err != nil {
		log.Printf("error closing Reader: %v", err)
	}
	return nil
}
