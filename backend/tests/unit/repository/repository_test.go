package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates an in-memory SQLite database with all necessary tables
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			price INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			name TEXT NOT NULL,
			description TEXT,
			thumbnail TEXT,
			images TEXT DEFAULT '[]',
			tags TEXT DEFAULT '[]',
			stock INTEGER DEFAULT 0,
			is_active INTEGER DEFAULT 1,
			status INTEGER DEFAULT 0
		);

		CREATE TABLE orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			total_amount INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			customer_phone TEXT,
			shipping_address TEXT,
			items TEXT DEFAULT '[]',
			payment_method INTEGER DEFAULT 0,
			status INTEGER DEFAULT 0,
			pay_os_order_code INTEGER UNIQUE
		);

		CREATE TABLE limited_drops (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_id INTEGER NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME,
			name TEXT NOT NULL,
			total_stock INTEGER NOT NULL DEFAULT 0,
			drop_size INTEGER NOT NULL DEFAULT 1,
			sold INTEGER NOT NULL DEFAULT 0,
			is_active INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE symbicodes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER,
			product_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			activated_at DATETIME,
			code BLOB NOT NULL UNIQUE,
			secret_key TEXT NOT NULL,
			activated_ip TEXT,
			is_activated INTEGER NOT NULL DEFAULT 0
		);
	`)
	require.NoError(t, err)

	return db
}

// =============================================================================
// PRODUCT REPOSITORY TESTS
// =============================================================================

func TestGetProductByID_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	_, err := db.Exec(`INSERT INTO products (id, name, description, thumbnail, price, stock, is_active, images, tags) VALUES (1, 'Test Watch', 'A test watch', '/img/watch.jpg', 500000, 10, 1, '[]', '[]')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO products (id, name, description, thumbnail, price, stock, is_active, images, tags) VALUES (2, 'Inactive Watch', 'Inactive desc', '/img/inactive.jpg', 300000, 5, 0, '[]', '[]')`)
	require.NoError(t, err)

	repo := repository.NewRepository(db)

	tests := []struct {
		name      string
		productID uint64
		wantErr   bool
		wantName  string
	}{
		{
			name:      "success - product exists",
			productID: 1,
			wantErr:   false,
			wantName:  "Test Watch",
		},
		{
			name:      "error - product not found",
			productID: 999,
			wantErr:   true,
		},
		{
			name:      "error - inactive product",
			productID: 2,
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			product, err := repo.GetProductByID(tc.productID)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, product)
			assert.Equal(t, tc.wantName, product.Name)
		})
	}
}

func TestGetAllProducts_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewRepository(db)

	tests := []struct {
		name      string
		setup     func()
		wantCount int
	}{
		{
			name:      "empty - no products",
			setup:     func() {},
			wantCount: 0,
		},
		{
			name: "success - returns only active products",
			setup: func() {
				db.Exec(`INSERT INTO products (name, description, thumbnail, price, is_active, images, tags) VALUES ('Active1', 'desc1', '/img/1.jpg', 100000, 1, '[]', '[]')`)
				db.Exec(`INSERT INTO products (name, description, thumbnail, price, is_active, images, tags) VALUES ('Active2', 'desc2', '/img/2.jpg', 200000, 1, '[]', '[]')`)
				db.Exec(`INSERT INTO products (name, description, thumbnail, price, is_active, images, tags) VALUES ('Inactive', 'desc3', '/img/3.jpg', 300000, 0, '[]', '[]')`)
			},
			wantCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset database
			db.Exec(`DELETE FROM products`)
			tc.setup()

			products, err := repo.GetAllProducts()

			assert.NoError(t, err)
			assert.Len(t, products, tc.wantCount)
		})
	}
}

// =============================================================================
// ORDER REPOSITORY TESTS
// =============================================================================

func TestCreateOrder_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewRepository(db)

	tests := []struct {
		name    string
		order   *models.Order
		wantErr bool
	}{
		{
			name: "success - create order with PayOS code",
			order: &models.Order{
				TotalAmount:     500000,
				CreatedAt:       time.Now(),
				CustomerPhone:   "0123456789",
				ShippingAddress: []byte(`{"name":"John"}`),
				Items:           []byte(`[{"product_id":1}]`),
				PaymentMethod:   1,
				Status:          0,
				PayOSOrderCode:  ptrInt64(12345),
			},
			wantErr: false,
		},
		{
			name: "success - create order without PayOS code",
			order: &models.Order{
				TotalAmount:     300000,
				CreatedAt:       time.Now(),
				CustomerPhone:   "0987654321",
				ShippingAddress: []byte(`{}`),
				Items:           []byte(`[]`),
				PaymentMethod:   0,
				Status:          0,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateOrder(tc.order)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, tc.order.ID)
		})
	}
}

func TestGetOrderByID_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	_, err := db.Exec(`INSERT INTO orders (id, total_amount, customer_phone, shipping_address, items, status) VALUES (1, 500000, '0123', '{}', '[]', 0)`)
	require.NoError(t, err)

	repo := repository.NewRepository(db)

	tests := []struct {
		name    string
		orderID uint64
		wantErr bool
	}{
		{
			name:    "success - order exists",
			orderID: 1,
			wantErr: false,
		},
		{
			name:    "error - order not found",
			orderID: 999,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			order, err := repo.GetOrderByID(tc.orderID)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, order)
			assert.Equal(t, tc.orderID, order.ID)
		})
	}
}

func TestGetOrdersByUserPhone_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	db.Exec(`INSERT INTO orders (total_amount, customer_phone, shipping_address, items, status) VALUES (100000, '0123', '{}', '[]', 0)`)
	db.Exec(`INSERT INTO orders (total_amount, customer_phone, shipping_address, items, status) VALUES (200000, '0123', '{}', '[]', 2)`)
	db.Exec(`INSERT INTO orders (total_amount, customer_phone, shipping_address, items, status) VALUES (300000, '9999', '{}', '[]', 0)`)

	repo := repository.NewRepository(db)

	tests := []struct {
		name      string
		phone     string
		wantCount int
	}{
		{
			name:      "success - user has multiple orders",
			phone:     "0123",
			wantCount: 2,
		},
		{
			name:      "success - user has one order",
			phone:     "9999",
			wantCount: 1,
		},
		{
			name:      "success - user has no orders",
			phone:     "0000",
			wantCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			orders, err := repo.GetOrdersByUserPhone(tc.phone)

			assert.NoError(t, err)
			assert.Len(t, orders, tc.wantCount)
		})
	}
}

// =============================================================================
// DROP REPOSITORY TESTS
// =============================================================================

func TestGetDropByID_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	now := time.Now()
	// Seed data
	db.Exec(`INSERT INTO limited_drops (id, product_id, start_time, name, total_stock, drop_size, sold, is_active) VALUES (1, 10, ?, 'Active Drop', 100, 50, 5, 1)`, now)
	db.Exec(`INSERT INTO limited_drops (id, product_id, start_time, name, total_stock, drop_size, sold, is_active) VALUES (2, 20, ?, 'Inactive Drop', 50, 25, 0, 0)`, now)

	repo := repository.NewRepository(db)

	tests := []struct {
		name    string
		dropID  uint64
		wantErr bool
	}{
		{
			name:    "success - active drop exists",
			dropID:  1,
			wantErr: false,
		},
		{
			name:    "error - inactive drop",
			dropID:  2,
			wantErr: true,
		},
		{
			name:    "error - drop not found",
			dropID:  999,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			drop, err := repo.GetDropByID(tc.dropID)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, drop)
		})
	}
}

func TestGetActiveDrops_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewRepository(db)
	now := time.Now()

	tests := []struct {
		name      string
		setup     func()
		wantCount int
	}{
		{
			name:      "empty - no drops",
			setup:     func() {},
			wantCount: 0,
		},
		{
			name: "success - returns only active drops",
			setup: func() {
				db.Exec(`INSERT INTO limited_drops (product_id, start_time, name, total_stock, is_active) VALUES (1, ?, 'Drop1', 10, 1)`, now)
				db.Exec(`INSERT INTO limited_drops (product_id, start_time, name, total_stock, is_active) VALUES (2, ?, 'Drop2', 20, 1)`, now)
				db.Exec(`INSERT INTO limited_drops (product_id, start_time, name, total_stock, is_active) VALUES (3, ?, 'Inactive', 30, 0)`, now)
			},
			wantCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec(`DELETE FROM limited_drops`)
			tc.setup()

			drops, err := repo.GetActiveDrops()

			assert.NoError(t, err)
			assert.Len(t, drops, tc.wantCount)
		})
	}
}

func TestIncrementSoldCount_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewRepository(db)
	now := time.Now()

	tests := []struct {
		name      string
		setup     func()
		dropID    uint64
		increment uint32
		wantErr   bool
	}{
		{
			name: "success - increment within stock",
			setup: func() {
				db.Exec(`DELETE FROM limited_drops`)
				db.Exec(`INSERT INTO limited_drops (id, product_id, start_time, name, total_stock, sold, is_active) VALUES (1, 10, ?, 'Drop', 100, 50, 1)`, now)
			},
			dropID:    1,
			increment: 10,
			wantErr:   false,
		},
		{
			name: "error - sold out (increment exceeds stock)",
			setup: func() {
				db.Exec(`DELETE FROM limited_drops`)
				db.Exec(`INSERT INTO limited_drops (id, product_id, start_time, name, total_stock, sold, is_active) VALUES (1, 10, ?, 'Drop', 100, 95, 1)`, now)
			},
			dropID:    1,
			increment: 10,
			wantErr:   true,
		},
		{
			name: "error - already sold out",
			setup: func() {
				db.Exec(`DELETE FROM limited_drops`)
				db.Exec(`INSERT INTO limited_drops (id, product_id, start_time, name, total_stock, sold, is_active) VALUES (1, 10, ?, 'Drop', 100, 100, 1)`, now)
			},
			dropID:    1,
			increment: 1,
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			err := repo.IncrementSoldCount(tc.dropID, tc.increment)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

// =============================================================================
// SYMBICODE REPOSITORY TESTS
// =============================================================================

func TestCreateSymbicode_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewRepository(db)

	tests := []struct {
		name      string
		symbicode *models.Symbicode
		wantErr   bool
	}{
		{
			name: "success - create symbicode",
			symbicode: &models.Symbicode{
				OrderID:   1,
				ProductID: 10,
				CreatedAt: time.Now(),
				Code:      []byte{0x01, 0x02, 0x03, 0x04},
				SecretKey: "secret123",
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateSymbicode(tc.symbicode)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestGetSymbicodeByCode_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	testCode := []byte{0xAB, 0xCD, 0xEF}
	db.Exec(`INSERT INTO symbicodes (order_id, product_id, code, secret_key, activated_ip, is_activated) VALUES (1, 10, ?, 'secret', '', 0)`, testCode)

	repo := repository.NewRepository(db)

	tests := []struct {
		name    string
		code    []byte
		wantErr bool
	}{
		{
			name:    "success - symbicode exists",
			code:    testCode,
			wantErr: false,
		},
		{
			name:    "error - symbicode not found",
			code:    []byte{0x00, 0x00},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			symbicode, err := repo.GetSymbicodeByCode(tc.code)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, symbicode)
		})
	}
}

func TestActivateSymbicode_TableDriven(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	db.Exec(`INSERT INTO symbicodes (id, order_id, product_id, code, secret_key, is_activated) VALUES (1, 1, 10, X'ABCD', 'secret', 0)`)

	repo := repository.NewRepository(db)

	tests := []struct {
		name        string
		symbicodeID uint64
		ip          string
		wantErr     bool
	}{
		{
			name:        "success - activate symbicode",
			symbicodeID: 1,
			ip:          "192.168.1.1",
			wantErr:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.ActivateSymbicode(tc.symbicodeID, tc.ip)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify activation
			var isActivated int
			db.QueryRow(`SELECT is_activated FROM symbicodes WHERE id = ?`, tc.symbicodeID).Scan(&isActivated)
			assert.Equal(t, 1, isActivated)
		})
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func ptrInt64(v int64) *int64 {
	return &v
}
