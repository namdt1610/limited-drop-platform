package service_test

import (
	"errors"
	"testing"
	"time"

	"ecommerce-backend/internal/integrations"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"

	"gorm.io/datatypes"
)

// =============================================================================
// DROP SERVICE TESTS
// =============================================================================

func TestGetActiveDrops_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name: "success - has active drops",
			setup: func(m *mockRepository) {
				m.activeDrops = []models.LimitedDrop{
					{ID: 1, Name: "Drop1"},
					{ID: 2, Name: "Drop2"},
				}
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "success - no active drops",
			setup:     func(m *mockRepository) {},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "error - repository error",
			setup: func(m *mockRepository) {
				m.getActiveDropsErr = errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			tc.setup(repo)
			srv := service.NewService(repo, nil, nil, nil)

			drops, err := srv.GetActiveDrops()

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if len(drops) != tc.wantCount {
				t.Fatalf("expected %d drops, got %d", tc.wantCount, len(drops))
			}
		})
	}
}

func TestGetDropStatus_TableDriven(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		dropID        uint64
		setup         func(*mockRepository)
		wantErr       bool
		wantAvailable uint32
	}{
		{
			name:   "success - drop with available stock",
			dropID: 1,
			setup: func(m *mockRepository) {
				m.drops[1] = &models.LimitedDrop{
					ID: 1, ProductID: 10, Name: "Test Drop",
					TotalStock: 100, Sold: 30, DropSize: 50,
					StartTime: now.Add(-time.Hour), IsActive: 1,
				}
				m.products[10] = &models.Product{ID: 10, Name: "Test Product", Price: 500000}
			},
			wantErr:       false,
			wantAvailable: 70,
		},
		{
			name:   "success - drop sold out",
			dropID: 2,
			setup: func(m *mockRepository) {
				m.drops[2] = &models.LimitedDrop{
					ID: 2, ProductID: 10, Name: "Sold Out Drop",
					TotalStock: 10, Sold: 10, DropSize: 10,
					StartTime: now.Add(-time.Hour), IsActive: 1,
				}
				m.products[10] = &models.Product{ID: 10, Name: "Test Product", Price: 500000}
			},
			wantErr:       false,
			wantAvailable: 0,
		},
		{
			name:    "error - drop not found",
			dropID:  999,
			setup:   func(m *mockRepository) {},
			wantErr: true,
		},
		{
			name:   "error - product not found",
			dropID: 1,
			setup: func(m *mockRepository) {
				m.drops[1] = &models.LimitedDrop{
					ID: 1, ProductID: 999, Name: "Test Drop",
					TotalStock: 100, Sold: 0,
					StartTime: now.Add(-time.Hour), IsActive: 1,
				}
				// No product with ID 999
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			tc.setup(repo)
			srv := service.NewService(repo, nil, nil, nil)

			status, err := srv.GetDropStatus(tc.dropID)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if status.Available != tc.wantAvailable {
				t.Fatalf("expected available %d, got %d", tc.wantAvailable, status.Available)
			}
		})
	}
}

func TestPurchaseDrop_TableDriven(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name           string
		dropID         uint64
		request        *service.PurchaseRequest
		setup          func(*mockRepository, *mockPaymentGateway)
		wantErr        string
		wantPaymentURL string
	}{
		{
			name:   "success - valid purchase with PayOS",
			dropID: 1,
			request: &service.PurchaseRequest{
				Quantity: 1, Name: "John", Phone: "0123",
				Email:    "john@test.com", Address: "123 St",
				Province: "HCM", District: "D1", Ward: "W1",
			},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				m.drops[1] = &models.LimitedDrop{
					ID:         1,
					ProductID:  10,
					TotalStock: 100,
					Sold:       0,
					DropSize:   50,
					StartTime:  now.Add(-time.Minute),
					IsActive:   1,
				}
				m.products[10] = &models.Product{ID: 10, Name: "Test", Price: 100000}
				pg.checkoutResponse = &integrations.PayOSCheckoutResponse{
					Data: struct {
						Bin           string `json:"bin"`
						AccountNumber string `json:"accountNumber"`
						AccountName   string `json:"accountName"`
						Amount        int64  `json:"amount"`
						Description   string `json:"description"`
						OrderCode     int64  `json:"orderCode"`
						Currency      string `json:"currency"`
						PaymentLinkID string `json:"paymentLinkId"`
						QRCode        string `json:"qrCode"`
						CheckoutURL   string `json:"checkoutUrl"`
					}{
						CheckoutURL: "https://payos.vn/checkout",
					},
				}
			},
			wantPaymentURL: "https://payos.vn/checkout",
		},
		{
			name:    "error - drop not active",
			dropID:  1,
			request: &service.PurchaseRequest{Quantity: 1, Name: "John", Phone: "0123"},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				m.drops[1] = &models.LimitedDrop{
					ID:        1,
					ProductID: 10,
					StartTime: now.Add(-time.Minute),
					IsActive:  0,
				}
			},
			wantErr: "limited drop is not active",
		},
		{
			name:    "error - drop not started",
			dropID:  1,
			request: &service.PurchaseRequest{Quantity: 1, Name: "John", Phone: "0123"},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				m.drops[1] = &models.LimitedDrop{
					ID:        1,
					ProductID: 10,
					StartTime: now.Add(time.Hour),
					IsActive:  1,
				}
			},
			wantErr: "limited drop has not started yet",
		},
		{
			name:    "error - drop has ended",
			dropID:  1,
			request: &service.PurchaseRequest{Quantity: 1, Name: "John", Phone: "0123"},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				endTime := now.Add(-time.Minute)
				m.drops[1] = &models.LimitedDrop{
					ID:        1,
					ProductID: 10,
					StartTime: now.Add(-2 * time.Hour),
					EndTime:   &endTime,
					IsActive:  1,
				}
			},
			wantErr: "limited drop has ended",
		},
		{
			name:    "error - sold out (total stock)",
			dropID:  1,
			request: &service.PurchaseRequest{Quantity: 1, Name: "John", Phone: "0123"},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				m.drops[1] = &models.LimitedDrop{
					ID:         1,
					ProductID:  10,
					TotalStock: 100,
					Sold:       100,
					StartTime:  now.Add(-time.Minute),
					IsActive:   1,
				}
			},
			wantErr: "limited drop is sold out",
		},
		{
			name:    "error - drop size limit reached",
			dropID:  1,
			request: &service.PurchaseRequest{Quantity: 1, Name: "John", Phone: "0123"},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				m.drops[1] = &models.LimitedDrop{
					ID:         1,
					ProductID:  10,
					TotalStock: 100,
					Sold:       50,
					DropSize:   50,
					StartTime:  now.Add(-time.Minute),
					IsActive:   1,
				}
			},
			wantErr: "limited drop size limit reached",
		},
		{
			name:    "error - drop not found",
			dropID:  999,
			request: &service.PurchaseRequest{Quantity: 1, Name: "John", Phone: "0123"},
			setup:   func(m *mockRepository, pg *mockPaymentGateway) {},
			wantErr: "drop not found",
		},
		{
			name:    "error - product not found",
			dropID:  1,
			request: &service.PurchaseRequest{Quantity: 1, Name: "John", Phone: "0123"},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				m.drops[1] = &models.LimitedDrop{
					ID: 1, ProductID: 999,
					TotalStock: 100, Sold: 0, DropSize: 50,
					StartTime: now.Add(-time.Minute), IsActive: 1,
				}
				// No product 999
			},
			wantErr: "product not found",
		},
		{
			name:    "error - payment gateway failure",
			dropID:  1,
			request: &service.PurchaseRequest{Quantity: 1, Name: "Test", Phone: "0123"},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				m.drops[1] = &models.LimitedDrop{ID: 1, ProductID: 10, TotalStock: 10, DropSize: 50, StartTime: now.Add(-time.Minute), IsActive: 1}
				m.products[10] = &models.Product{ID: 10, Name: "Test", Price: 100000}
				pg.checkoutErr = errors.New("payos error")
			},
			wantErr: "failed to create PayOS checkout",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			pg := newMockPaymentGateway()
			
			if tc.setup != nil {
				tc.setup(repo, pg)
			}
			
			srv := service.NewService(repo, pg, nil, nil)

			result, err := srv.PurchaseDrop(tc.dropID, tc.request)

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error '%s', got nil", tc.wantErr)
				}
				if err.Error() != tc.wantErr && !contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error '%s', got '%s'", tc.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if result == nil {
				t.Fatal("expected result, got nil")
			}
			if tc.wantPaymentURL != "" && result.PaymentURL != tc.wantPaymentURL {
				t.Errorf("expected payment URL '%s', got '%s'", tc.wantPaymentURL, result.PaymentURL)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}


func TestProcessSuccessfulDropPayment_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		orderCode int64
		dropID    uint64
		setup     func(*mockRepository, *mockEmailSender, *mockSheetSubmitter)
		wantErr   bool
		wantSold  uint32 // expected sold count after processing
	}{
		{
			name:      "success - winner (first payment)",
			orderCode: 12345,
			dropID:    1,
			setup: func(m *mockRepository, e *mockEmailSender, s *mockSheetSubmitter) {
				m.drops[1] = &models.LimitedDrop{
					ID: 1, ProductID: 10,
					TotalStock: 100, Sold: 0, DropSize: 50,
				}
				m.products[10] = &models.Product{ID: 10, Price: 100000}

				order := &models.Order{
					ID:              100,
					Status:          models.OrderPending,
					Items:           datatypes.JSON(`[{"drop_id":1,"quantity":1,"product_id":10}]`),
					ShippingAddress: datatypes.JSON(`{"name":"Winner","email":"winner@test.com","phone":"0123456789"}`),
					CustomerPhone:   "0123456789",
				}
				m.orderByPayOS[12345] = order
				m.orders[100] = order
			},
			wantErr:  false,
			wantSold: 1,
		},
		{
			name:      "idempotency - already processed",
			orderCode: 12345,
			dropID:    1,
			setup: func(m *mockRepository, e *mockEmailSender, s *mockSheetSubmitter) {
				m.drops[1] = &models.LimitedDrop{
					ID: 1, TotalStock: 100, Sold: 1,
				}
				m.orderByPayOS[12345] = &models.Order{
					ID:     100,
					Status: models.OrderPaid,
					Items:  datatypes.JSON(`[{"drop_id":1,"quantity":1,"product_id":10}]`),
				}
			},
			wantErr:  false,
			wantSold: 1,
		},
		{
			name:      "loser - sold out during race",
			orderCode: 99999,
			dropID:    1,
			setup: func(m *mockRepository, e *mockEmailSender, s *mockSheetSubmitter) {
				m.drops[1] = &models.LimitedDrop{
					ID: 1, ProductID: 10,
					TotalStock: 10, Sold: 10, // Already sold out
				}
				m.allowIncrement = false // Force ErrSoldOut

				m.orderByPayOS[99999] = &models.Order{
					ID:              200,
					Status:          models.OrderPending,
					Items:           datatypes.JSON(`[{"drop_id":1,"quantity":1,"product_id":10}]`),
					ShippingAddress: datatypes.JSON(`{"name":"Loser","email":"loser@test.com","phone":"0123456789"}`),
					CustomerPhone:   "0123456789",
				}
			},
			wantErr:  false,
			wantSold: 10,
		},
		{
			name:      "error - increment error (not sold out)",
			orderCode: 77777,
			dropID:    1,
			setup: func(m *mockRepository, e *mockEmailSender, s *mockSheetSubmitter) {
				m.drops[1] = &models.LimitedDrop{
					ID: 1, ProductID: 10,
					TotalStock: 100, Sold: 0,
				}
				m.incrementErr = errors.New("database error")
				m.orderByPayOS[77777] = &models.Order{
					ID:     300,
					Status: models.OrderPending,
					Items:  datatypes.JSON(`[{"drop_id":1,"quantity":1,"product_id":10}]`),
				}
			},
			wantErr:  true,
			wantSold: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			email := newMockEmailSender()
			sheets := newMockSheetSubmitter()
			tc.setup(repo, email, sheets)
			srv := service.NewService(repo, nil, email, sheets)

			err := srv.ProcessSuccessfulDropPayment(tc.orderCode)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}

			// Verify sold count
			if drop, ok := repo.drops[tc.dropID]; ok {
				if drop.Sold != tc.wantSold {
					t.Fatalf("expected sold=%d, got sold=%d", tc.wantSold, drop.Sold)
				}
			}
		})
	}
}

// =============================================================================
// EDGE CASE TESTS
// =============================================================================

func TestPurchaseDrop_EdgeCases(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		dropID  uint64
		request *service.PurchaseRequest
		setup   func(*mockRepository, *mockPaymentGateway)
		wantErr string
	}{
		{
			name:   "edge - exactly at start time",
			dropID: 1,
			request: &service.PurchaseRequest{
				Quantity: 1, Name: "John", Phone: "0123",
			},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				m.drops[1] = &models.LimitedDrop{
					ID: 1, ProductID: 10,
					TotalStock: 100, Sold: 0, DropSize: 50,
					StartTime: now, // Exactly now
					IsActive:  1,
				}
				m.products[10] = &models.Product{ID: 10, Name: "Test", Price: 100000}
				pg.checkoutResponse = &integrations.PayOSCheckoutResponse{
					Data: struct {
						Bin           string `json:"bin"`
						AccountNumber string `json:"accountNumber"`
						AccountName   string `json:"accountName"`
						Amount        int64  `json:"amount"`
						Description   string `json:"description"`
						OrderCode     int64  `json:"orderCode"`
						Currency      string `json:"currency"`
						PaymentLinkID string `json:"paymentLinkId"`
						QRCode        string `json:"qrCode"`
						CheckoutURL   string `json:"checkoutUrl"`
					}{
						CheckoutURL: "https://payos.vn/checkout",
					},
				}
			},
			wantErr: "", // Should succeed
		},
		{
			name:   "edge - sold equals drop size minus one",
			dropID: 1,
			request: &service.PurchaseRequest{
				Quantity: 1, Name: "John", Phone: "0123",
			},
			setup: func(m *mockRepository, pg *mockPaymentGateway) {
				m.drops[1] = &models.LimitedDrop{
					ID: 1, ProductID: 10,
					TotalStock: 100, Sold: 49, DropSize: 50,
					StartTime: now.Add(-time.Minute), IsActive: 1,
				}
				m.products[10] = &models.Product{ID: 10, Name: "Test", Price: 100000}
				pg.checkoutResponse = &integrations.PayOSCheckoutResponse{
					Data: struct {
						Bin           string `json:"bin"`
						AccountNumber string `json:"accountNumber"`
						AccountName   string `json:"accountName"`
						Amount        int64  `json:"amount"`
						Description   string `json:"description"`
						OrderCode     int64  `json:"orderCode"`
						Currency      string `json:"currency"`
						PaymentLinkID string `json:"paymentLinkId"`
						QRCode        string `json:"qrCode"`
						CheckoutURL   string `json:"checkoutUrl"`
					}{
						CheckoutURL: "https://payos.vn/checkout",
					},
				}
			},
			wantErr: "", // Should succeed
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			pg := newMockPaymentGateway()
			tc.setup(repo, pg)
			srv := service.NewService(repo, pg, nil, nil)

			result, err := srv.PurchaseDrop(tc.dropID, tc.request)

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error '%s', got nil", tc.wantErr)
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected error '%s', got '%s'", tc.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if result == nil {
				t.Fatal("expected result, got nil")
			}
		})
	}
}
