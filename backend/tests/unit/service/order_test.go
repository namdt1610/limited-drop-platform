package service_test

import (
	"errors"
	"testing"
	"time"

	"ecommerce-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrder_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		customerPhone  string
		shippingAddr   []byte
		items          []byte
		paymentMethod  uint8
		payOSOrderCode *int64
		mockError      error
		wantErr        bool
		wantTotal      uint64
	}{
		{
			name:          "success",
			customerPhone: "0909",
			shippingAddr:  []byte(`{"address":"123"}`),
			items:         []byte(`[{"product_name":"A","price":100,"quantity":2}]`),
			paymentMethod: 1,
			mockError:     nil,
			wantErr:       false,
			wantTotal:     200,
		},
		{
			name:          "error - db failed",
			customerPhone: "0909",
			shippingAddr:  []byte(`{}`),
			items:         []byte(`[]`),
			mockError:     errors.New("db error"),
			wantErr:       true,
			wantTotal:     0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, m := setup()
			if tc.mockError != nil {
				m.createOrderErr = tc.mockError
			}

			order, err := s.CreateOrder(tc.customerPhone, tc.shippingAddr, tc.items, tc.paymentMethod, tc.payOSOrderCode)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantTotal, order.TotalAmount)
			}
		})
	}
}

func TestGetOrderByID_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		orderID    uint64
		mockReturn *models.Order
		mockError  error
		wantErr    bool
	}{
		{
			name:       "success",
			orderID:    1,
			mockReturn: &models.Order{ID: 1},
			mockError:  nil,
			wantErr:    false,
		},
		{
			name:       "error - not found",
			orderID:    99,
			mockReturn: nil,
			mockError:  errors.New("not found"),
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, m := setup()
			if tc.mockError != nil {
				m.getOrderErr = tc.mockError
			} else if tc.mockReturn != nil {
				m.orders[tc.orderID] = tc.mockReturn
			}

			order, err := s.GetOrderByID(tc.orderID)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockReturn, order)
			}
		})
	}
}

func TestGetOrdersByUserPhone_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		phone      string
		mockReturn []models.Order
		mockError  error
		wantErr    bool
	}{
		{
			name:  "success",
			phone: "0909",
			mockReturn: []models.Order{
				{ID: 1, CreatedAt: time.Now()},
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name:       "error - db error",
			phone:      "0909",
			mockReturn: nil,
			mockError:  errors.New("db error"),
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, m := setup()
			if tc.mockError != nil {
				m.getOrdersErr = tc.mockError
			} else {
				m.ordersByPhone[tc.phone] = tc.mockReturn
			}

			orders, err := s.GetOrdersByUserPhone(tc.phone)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockReturn, orders)
			}
		})
	}
}
