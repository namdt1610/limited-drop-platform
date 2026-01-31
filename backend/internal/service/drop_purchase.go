package service

import (
	"ecommerce-backend/internal/integrations"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// PurchaseDrop handles the business logic for purchasing drop items
func (s *service) PurchaseDrop(dropID uint64, req *PurchaseRequest) (*PurchaseResult, error) {
	// Get the drop
	drop, err := s.repo.GetDropByID(dropID)
	if err != nil {
		return nil, err
	}

	// Check if drop is active
	if drop.IsActive != 1 {
		return nil, errors.New("limited drop is not active")
	}

	// Check if drop is still running
	now := time.Now()
	if now.Before(drop.StartTime) {
		return nil, errors.New("limited drop has not started yet")
	}
	if drop.EndTime != nil && now.After(*drop.EndTime) {
		return nil, errors.New("limited drop has ended")
	}

	// Check if stock is available
	if drop.Sold >= drop.TotalStock {
		return nil, errors.New("limited drop is sold out")
	}

	// Check drop size limit
	if drop.Sold >= drop.DropSize {
		return nil, errors.New("limited drop size limit reached")
	}

	// Get the product for pricing
	product, err := s.GetProduct(drop.ProductID)
	if err != nil {
		return nil, err
	}

	// Create PayOS checkout
	// Use UnixNano to avoid collisions in high traffic
	orderCode := time.Now().UnixNano()
	amount := product.Price * uint64(req.Quantity)

	// Create shipping address JSON
	shippingAddress := map[string]interface{}{
		"name":     req.Name,
		"phone":    req.Phone,
		"email":    req.Email,
		"address":  req.Address,
		"province": req.Province,
		"district": req.District,
		"ward":     req.Ward,
	}
	shippingJSON, _ := json.Marshal(shippingAddress)

	// Create items JSON
	items := []map[string]interface{}{
		{
			"product_id": dropID,
			"drop_id":    dropID,
			"name":       product.Name,
			"price":      product.Price,
			"quantity":   req.Quantity,
		},
	}
	itemsJSON, _ := json.Marshal(items)

	// Create order in database FIRST with PENDING payment status (status = 1)
	// This ensures that if payment is successful, we definitely have the order record.
	// Pass PayOSOrderCode to CreateOrder to link the transaction
	_, err = s.CreateOrder(req.Phone, shippingJSON, itemsJSON, 1, &orderCode) // 1 = PayOS payment method
	if err != nil {
		return nil, fmt.Errorf("failed to create local order: %w", err)
	}

	// Get frontend URL from environment, default to localhost:3000
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	payosReq := integrations.PayOSCheckoutRequest{
		OrderCode:   orderCode,
		Amount:      int64(amount),
		Description: fmt.Sprintf("Drop %d", dropID),
		ReturnURL:   frontendURL + "/#payment-success",
		CancelURL:   frontendURL + "/#payment-cancel",
		Items: []integrations.PayOSItem{
			{
				Name:     product.Name,
				Quantity: req.Quantity,
				Price:    int64(product.Price),
			},
		},
	}

	checkout, err := s.payment.CreateCheckout(payosReq)
	if err != nil {
		// If checkout creation fails, the order remains as PENDING (Abandoned Cart)
		return nil, fmt.Errorf("failed to create PayOS checkout: %w", err)
	}

	return &PurchaseResult{
		Message:    "Đơn hàng đã được tạo!",
		PaymentURL: checkout.Data.CheckoutURL,
		OrderCode:  orderCode,
	}, nil
}
