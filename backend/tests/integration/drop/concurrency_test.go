package drop_test

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"ecommerce-backend/internal/integrations"
	"ecommerce-backend/internal/repository"
	servicepkg "ecommerce-backend/internal/service"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	// Force single connection to avoid SQLite in-memory per-connection issue
	db.SetMaxOpenConns(1)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	// Create minimal limited_drops table used by tests
	schema := `CREATE TABLE limited_drops (
		id INTEGER PRIMARY KEY,
		product_id INTEGER,
		start_time DATETIME,
		end_time DATETIME,
		name TEXT,
		total_stock INTEGER,
		drop_size INTEGER,
		sold INTEGER DEFAULT 0,
		is_active INTEGER
	);`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	return db
}

func insertDrop(t *testing.T, db *sql.DB, id int64, start time.Time, dropSize, totalStock, sold, isActive int) {
	query := `INSERT INTO limited_drops (id, product_id, start_time, end_time, name, total_stock, drop_size, sold, is_active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := db.Exec(query, id, 1, start, nil, fmt.Sprintf("drop-%d", id), totalStock, dropSize, sold, isActive); err != nil {
		t.Fatalf("insert drop: %v", err)
	}
}

func TestPurchaseDrop_Concurrency(t *testing.T) {
	tests := []struct {
		name     string
		attempts int
		dropSize int
		stock    int
	}{
		{"concurrency_high_stock", 100, 100, 100},
		{"concurrency_more_attempts_than_stock", 200, 200, 50},
		{"concurrency_drop_size_limited", 100, 10, 100},
		{"concurrency_small", 10, 5, 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := setupDB(t)
			defer db.Close()

			now := time.Now().Add(-time.Minute)
			insertDrop(t, db, 1, now, tc.dropSize, tc.stock, 0, 1)

			repo := repository.NewRepository(db)

			// Mocks
			payment := &MockPaymentGateway{}
			email := &MockEmailSender{}
			sheets := &MockSheetSubmitter{}

			svc := servicepkg.NewService(repo, payment, email, sheets)

			var wg sync.WaitGroup
			wg.Add(tc.attempts)
			var mu sync.Mutex
			success := 0

			for i := 0; i < tc.attempts; i++ {
				go func(i int) {
					defer wg.Done()
					req := &servicepkg.PurchaseRequest{
						Quantity: 1,
						Name:     fmt.Sprintf("User %d", i),
						Phone:    fmt.Sprintf("phone-%d", i),
						Email:    fmt.Sprintf("user%d@example.com", i),
						Address:  "Test Address",
						Province: "Test Province",
						District: "Test District",
						Ward:     "Test Ward",
					}
					_, err := svc.PurchaseDrop(1, req)
					if err == nil {
						mu.Lock()
						success++
						mu.Unlock()
					}
				}(i)
			}
			wg.Wait()

			// read final sold
			var finalSold int
			if err := db.QueryRow("SELECT sold FROM limited_drops WHERE id = ?", 1).Scan(&finalSold); err != nil {
				t.Fatalf("query final sold: %v", err)
			}

			maxAllowed := tc.stock
			if tc.dropSize < maxAllowed {
				maxAllowed = tc.dropSize
			}

			if finalSold > maxAllowed {
				t.Fatalf("oversold: finalSold=%d maxAllowed=%d", finalSold, maxAllowed)
			}
			if success != finalSold {
				t.Fatalf("mismatch success count (%d) and finalSold (%d)", success, finalSold)
			}
		})
	}
}

// Mocks

type MockPaymentGateway struct{}

func (m *MockPaymentGateway) CreateCheckout(req integrations.PayOSCheckoutRequest) (*integrations.PayOSCheckoutResponse, error) {
	resp := &integrations.PayOSCheckoutResponse{}
	resp.Data.CheckoutURL = "http://mock-checkout-url"
	return resp, nil
}

func (m *MockPaymentGateway) VerifyPayment(orderCode int64) (*integrations.PayOSVerifyResponse, error) {
	return nil, nil
}

func (m *MockPaymentGateway) RefundPayment(orderCode int64, reason string) error {
	return nil
}

func (m *MockPaymentGateway) CancelPayment(orderCode int64) error {
	return nil
}

func (m *MockPaymentGateway) GenerateSignature(data string) string {
	return "mock-signature"
}

type MockEmailSender struct{}

func (m *MockEmailSender) SendOrderConfirmation(email, orderNumber string, amount float64) error {
	return nil
}

func (m *MockEmailSender) SendSymbioteReceipt(email, phone, status, elapsed string) error {
	return nil
}

func (m *MockEmailSender) SendOrderDetails(email string, order interface{}) error {
	return nil
}

type MockSheetSubmitter struct{}

func (m *MockSheetSubmitter) SubmitOrder(name, phone, email, address, notes string, amount float64, timestamp interface{}) error {
	return nil
}
