package service_test

import (
	"errors"
	"testing"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"
)

// =============================================================================
// ORDER SERVICE TESTS
// =============================================================================

func TestCreateOrder_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		phone         string
		shippingAddr  []byte
		items         []byte
		paymentMethod uint8
		setup         func(*mockRepository)
		wantErr       bool
	}{
		{
			name:          "success - valid order",
			phone:         "0123456789",
			shippingAddr:  []byte(`{"address":"123 Main St"}`),
			items:         []byte(`[{"price":100000,"quantity":2}]`),
			paymentMethod: 1,
			setup:         func(m *mockRepository) {},
			wantErr:       false,
		},
		{
			name:          "success - empty items calculates zero total",
			phone:         "0123456789",
			shippingAddr:  []byte(`{}`),
			items:         []byte(`[]`),
			paymentMethod: 1,
			setup:         func(m *mockRepository) {},
			wantErr:       false,
		},
		{
			name:          "success - invalid items JSON still creates order",
			phone:         "0123456789",
			shippingAddr:  []byte(`{}`),
			items:         []byte(`invalid json`),
			paymentMethod: 1,
			setup:         func(m *mockRepository) {},
			wantErr:       false,
		},
		{
			name:          "error - repository create error",
			phone:         "0123456789",
			shippingAddr:  []byte(`{}`),
			items:         []byte(`[]`),
			paymentMethod: 1,
			setup: func(m *mockRepository) {
				m.createOrderErr = errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			tc.setup(repo)
			srv := service.NewService(repo)

			// Pass nil for payOSOrderCode as per new signature
			order, err := srv.CreateOrder(tc.phone, tc.shippingAddr, tc.items, tc.paymentMethod, nil)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if order == nil {
				t.Fatal("expected order, got nil")
			}
		})
	}
}

func TestGetOrderByID_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		orderID uint64
		setup   func(*mockRepository)
		wantErr bool
	}{
		{
			name:    "success - order exists",
			orderID: 1,
			setup: func(m *mockRepository) {
				m.orders[1] = &models.Order{ID: 1, CustomerPhone: "0123"}
			},
			wantErr: false,
		},
		{
			name:    "error - order not found",
			orderID: 999,
			setup:   func(m *mockRepository) {},
			wantErr: true,
		},
		{
			name:    "error - repository error",
			orderID: 1,
			setup: func(m *mockRepository) {
				m.getOrderErr = errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			tc.setup(repo)
			srv := service.NewService(repo)

			order, err := srv.GetOrderByID(tc.orderID)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if order == nil {
				t.Fatal("expected order, got nil")
			}
		})
	}
}

func TestGetOrdersByUserPhone_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		phone     string
		setup     func(*mockRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - user has orders",
			phone: "0123456789",
			setup: func(m *mockRepository) {
				m.ordersByPhone["0123456789"] = []models.Order{
					{ID: 1}, {ID: 2}, {ID: 3},
				}
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "success - user has no orders",
			phone:     "0000000000",
			setup:     func(m *mockRepository) {},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:  "error - repository error",
			phone: "0123456789",
			setup: func(m *mockRepository) {
				m.getOrdersErr = errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			tc.setup(repo)
			srv := service.NewService(repo)

			orders, err := srv.GetOrdersByUserPhone(tc.phone)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if len(orders) != tc.wantCount {
				t.Fatalf("expected %d orders, got %d", tc.wantCount, len(orders))
			}
		})
	}
}
