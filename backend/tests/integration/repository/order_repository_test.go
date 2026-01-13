package repository_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"ecommerce-backend/internal/models"
	repoPkg "ecommerce-backend/internal/repository"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	db.SetMaxOpenConns(1)

	// create orders table compatible with repository expectations
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		total_amount INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		customer_phone TEXT,
		shipping_address TEXT,
		items TEXT,
		payment_method INTEGER DEFAULT 0,
		status INTEGER DEFAULT 0,
		payos_order_code INTEGER UNIQUE
	)`)
	if err != nil {
		t.Fatalf("failed to create orders table: %v", err)
	}
	return db
}

func TestRepository_CreateAndGetOrder_TableDriven(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	r := repoPkg.NewRepository(db)

	tests := []struct {
		name    string
		order   *models.Order
		wantErr bool
	}{
		{
			name: "create_and_get_success",
			order: &models.Order{
				TotalAmount:     12345,
				CustomerPhone:   "012333",
				ShippingAddress: []byte(`{"addr":"a"}`),
				Items:           []byte(`[{}]`),
				PaymentMethod:   1,
				Status:          models.OrderPending,
			},
			wantErr: false,
		},
		{
			name: "get_by_id_no_rows",
			order: &models.Order{
				CustomerPhone: "no-op",
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "create_and_get_success" {
				// create
				err := r.CreateOrder(tc.order)
				if err != nil {
					t.Fatalf("CreateOrder failed: %v", err)
				}

				// fetch by phone
				orders, err := r.GetOrdersByUserPhone(tc.order.CustomerPhone)
				if err != nil {
					t.Fatalf("GetOrdersByUserPhone failed: %v", err)
				}
				if len(orders) != 1 {
					t.Fatalf("expected 1 order, got %d", len(orders))
				}

				got := orders[0]
				if got.CustomerPhone != tc.order.CustomerPhone {
					t.Fatalf("expected phone %s, got %s", tc.order.CustomerPhone, got.CustomerPhone)
				}

				// get by id
				fetched, err := r.GetOrderByID(got.ID)
				if err != nil {
					t.Fatalf("GetOrderByID failed: %v", err)
				}
				if fetched.CustomerPhone != tc.order.CustomerPhone {
					t.Fatalf("expected phone %s, got %s", tc.order.CustomerPhone, fetched.CustomerPhone)
				}
				return
			}

			// for the "get_by_id_no_rows" case: expect no rows when querying unknown id
			_, err := r.GetOrderByID(9999)
			if err == nil {
				t.Fatalf("expected error for missing order, got nil")
			}
		})
	}
}
