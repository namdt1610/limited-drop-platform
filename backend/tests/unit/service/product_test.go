package service_test

import (
	"errors"
	"testing"

	"ecommerce-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestGetProduct_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		productID  uint64
		mockReturn *models.Product
		mockError  error
		wantErr    bool
	}{
		{
			name:       "success",
			productID:  1,
			mockReturn: &models.Product{ID: 1, Name: "Test Product"},
			mockError:  nil,
			wantErr:    false,
		},
		{
			name:       "error - not found",
			productID:  99,
			mockReturn: nil,
			mockError:  errors.New("not found"),
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, m := setup()
			
			// Configure Fake
			if tc.mockError != nil {
				m.productErr = tc.mockError
			} else if tc.mockReturn != nil {
				m.products[tc.productID] = tc.mockReturn
			}

			product, err := s.GetProduct(tc.productID)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockReturn, product)
			}
		})
	}
}

func TestListProducts_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*mockRepository)
		wantErr      bool
		wantCount    int
	}{
		{
			name: "success - returns active only",
			setupMock: func(m *mockRepository) {
				m.products[1] = &models.Product{ID: 1, Name: "Active", IsActive: 1}
				m.products[2] = &models.Product{ID: 2, Name: "Inactive", IsActive: 0}
				m.products[3] = &models.Product{ID: 3, Name: "Active 2", IsActive: 1}
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "error - db failed",
			setupMock: func(m *mockRepository) {
				m.allProductsErr = errors.New("db error")
			},
			wantErr:    true,
			wantCount:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, m := setup()
			tc.setupMock(m)

			products, err := s.ListProducts()

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, products, tc.wantCount)
			}
		})
	}
}
