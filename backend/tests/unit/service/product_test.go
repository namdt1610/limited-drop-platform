package service_test

import (
	"errors"
	"testing"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"
)

// =============================================================================
// PRODUCT SERVICE TESTS
// =============================================================================

func TestGetProduct_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		productID  uint64
		setup      func(*mockRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:      "success - product exists",
			productID: 1,
			setup: func(m *mockRepository) {
				m.products[1] = &models.Product{ID: 1, Name: "Test", Price: 100000}
			},
			wantErr: false,
		},
		{
			name:      "error - product not found",
			productID: 999,
			setup: func(m *mockRepository) {
				// No products
			},
			wantErr:    true,
			wantErrMsg: "product not found",
		},
		{
			name:      "error - repository error",
			productID: 1,
			setup: func(m *mockRepository) {
				m.productErr = errors.New("database connection error")
			},
			wantErr:    true,
			wantErrMsg: "database connection error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			tc.setup(repo)
			srv := service.NewService(repo)

			product, err := srv.GetProduct(tc.productID)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error '%s', got nil", tc.wantErrMsg)
				}
				if err.Error() != tc.wantErrMsg {
					t.Fatalf("expected error '%s', got '%s'", tc.wantErrMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if product == nil {
				t.Fatal("expected product, got nil")
			}
		})
	}
}

func TestListProducts_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name: "success - returns only active products",
			setup: func(m *mockRepository) {
				m.products[1] = &models.Product{ID: 1, Name: "Active", IsActive: 1}
				m.products[2] = &models.Product{ID: 2, Name: "Inactive", IsActive: 0}
				m.products[3] = &models.Product{ID: 3, Name: "Active2", IsActive: 1}
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "success - no active products",
			setup: func(m *mockRepository) {
				m.products[1] = &models.Product{ID: 1, Name: "Inactive", IsActive: 0}
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "success - empty product list",
			setup: func(m *mockRepository) {
				// No products
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "error - repository error",
			setup: func(m *mockRepository) {
				m.allProductsErr = errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			tc.setup(repo)
			srv := service.NewService(repo)

			products, err := srv.ListProducts()

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if len(products) != tc.wantCount {
				t.Fatalf("expected %d products, got %d", tc.wantCount, len(products))
			}
		})
	}
}
