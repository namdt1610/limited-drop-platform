package service

import (
	"ecommerce-backend/internal/models"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// CreateOrder creates a new order with business logic validation, optional payOSOrderCode
func (s *service) CreateOrder(customerPhone string, shippingAddress []byte, items []byte, paymentMethod uint8, payOSOrderCode *int64) (*models.Order, error) {
	// Calculate total amount from items
	var totalAmount uint64
	var itemsData []map[string]interface{}
	if err := json.Unmarshal(items, &itemsData); err == nil {
		for _, item := range itemsData {
			if price, ok := item["price"].(float64); ok {
				if quantity, ok := item["quantity"].(float64); ok {
					totalAmount += uint64(price * quantity)
				}
			}
		}
	}

	order := &models.Order{
		CustomerPhone:   customerPhone,
		ShippingAddress: datatypes.JSON(shippingAddress),
		Items:           datatypes.JSON(items),
		PaymentMethod:   paymentMethod,
		Status:          models.OrderPending,
		TotalAmount:     totalAmount,
		CreatedAt:       time.Now(),
		PayOSOrderCode:  payOSOrderCode,
	}

	err := s.repo.CreateOrder(order)
	if err != nil {
		return order, err
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID for tracking purposes
func (s *service) GetOrderByID(id uint64) (*models.Order, error) {
	return s.repo.GetOrderByID(id)
}

// GetOrdersByUserPhone retrieves all orders for a user for order history/tracking
func (s *service) GetOrdersByUserPhone(phone string) ([]models.Order, error) {
	return s.repo.GetOrdersByUserPhone(phone)
}
