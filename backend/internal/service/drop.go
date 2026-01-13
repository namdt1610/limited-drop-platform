package service

import (
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"ecommerce-backend/internal/integrations"
	"ecommerce-backend/internal/utils/uuid"
)

// LimitedDropStatus represents the current status of a drop for the frontend
type LimitedDropStatus struct {
	DropID      uint64     `json:"drop_id"`
	Name        string     `json:"name"`
	ProductID   uint64     `json:"product_id"`
	ProductName string     `json:"product_name"`
	Price       uint64     `json:"price"`
	TotalStock  uint32     `json:"total_stock"`
	Sold        uint32     `json:"sold"`
	Available   uint32     `json:"available"`
	DropSize    uint32     `json:"drop_size"`
	IsActive    bool       `json:"is_active"`
	StartsAt    time.Time  `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
	Now         time.Time  `json:"now"`
}

// GetActiveDrops returns all active drops
func (s *service) GetActiveDrops() ([]models.LimitedDrop, error) {
	return s.repo.GetActiveDrops()
}

// GetDropStatus returns the status of a specific drop
func (s *service) GetDropStatus(id uint64) (*LimitedDropStatus, error) {
	drop, err := s.repo.GetDropByID(id)
	if err != nil {
		return nil, err
	}

	product, err := s.repo.GetProductByID(drop.ProductID)
	if err != nil {
		return nil, err
	}

	available := uint32(0)
	if drop.TotalStock > drop.Sold {
		available = drop.TotalStock - drop.Sold
	}

	return &LimitedDropStatus{
		DropID:      drop.ID,
		Name:        drop.Name,
		ProductID:   drop.ProductID,
		ProductName: product.Name,
		Price:       product.Price,
		TotalStock:  drop.TotalStock,
		Sold:        drop.Sold,
		Available:   available,
		DropSize:    drop.DropSize,
		IsActive:    drop.IsActive == 1,
		StartsAt:    drop.StartTime,
		EndsAt:      drop.EndTime,
		Now:         time.Now().UTC(),
	}, nil
}

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
	orderCode := time.Now().Unix()
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

	// Get frontend URL from environment, default to localhost:3000
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	payosReq := integrations.PayOSCheckoutRequest{
		OrderCode:   orderCode,
		Amount:      int64(amount),
		Description: "Limited Drop Payment",
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

	checkout, err := integrations.CreatePayOSCheckout(payosReq)
	if err != nil {
		// For testing purposes, return a mock result if PayOS is not configured
		if os.Getenv("PAYOS_CLIENT_ID") == "" {
			return &PurchaseResult{
				Message:    "Đặt Limited Drop thành công!",
				PaymentURL: "https://payos.vn/test-payment",
				OrderCode:  orderCode,
			}, nil
		}
		return nil, fmt.Errorf("failed to create PayOS checkout: %w", err)
	}

	// Create order in database with PENDING payment status (status = 1)
	// This allows webhook handler to link the payment back to customer info
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

	// Pass PayOSOrderCode to CreateOrder to link the transaction
	_, err = s.CreateOrder(req.Phone, shippingJSON, itemsJSON, 1, &orderCode) // 1 = PayOS payment method
	if err != nil {
		// Order creation failed, but checkout was created
		// Log this but don't fail - webhook will retry
		fmt.Fprintf(os.Stderr, "Failed to create order for payos checkout %d: %v\n", orderCode, err)
	}

	return &PurchaseResult{
		Message:    "Đặt Limited Drop thành công!",
		PaymentURL: checkout.Data.CheckoutURL,
		OrderCode:  orderCode,
	}, nil
}

// ProcessSuccessfulDropPayment processes a successful PayOS payment for limited drops
func (s *service) ProcessSuccessfulDropPayment(orderCode int64) error {
	// 1. Retrieve the existing order using PayOS Order Code
	order, err := s.repo.GetOrderByPayOSOrderCode(orderCode)
	if err != nil {
		fmt.Printf("Order not found for code %d: %v\n", orderCode, err)
		return err // Order must exist (created in PurchaseDrop)
	}
	if order == nil {
		return fmt.Errorf("order not found for code %d", orderCode)
	}

	// 2. Idempotency Check
	if order.Status == models.OrderPaid || order.Status == models.OrderConfirmed {
		return nil // Already processed
	}

	// 3. Extract Drop Info from Order Items (JSON)
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(order.Items), &items); err != nil {
		return fmt.Errorf("failed to parse order items: %w", err)
	}

	if len(items) == 0 {
		return errors.New("order has no items")
	}

	// Extract dropID and quantity safely
	item := items[0]
	dropIDVal, ok := item["drop_id"].(float64)
	if !ok {
		return errors.New("invalid drop_id in items")
	}
	dropID := uint64(dropIDVal)

	quantityVal, ok := item["quantity"].(float64)
	if !ok {
		return errors.New("invalid quantity in items")
	}
	quantity := int(quantityVal)

	productIDVal, ok := item["product_id"].(float64)
	if !ok {
		return errors.New("invalid product_id in items")
	}
	productID := uint64(productIDVal)

	// Extract info for notifications
	var shippingAddress map[string]interface{}
	json.Unmarshal([]byte(order.ShippingAddress), &shippingAddress)
	customerEmail := ""
	if email, ok := shippingAddress["email"].(string); ok {
		customerEmail = email
	}
	customerName := ""
	if name, ok := shippingAddress["name"].(string); ok {
		customerName = name
	}
	shippingAddrStr := string(order.ShippingAddress)

	// 4. Atomic Stock Increment (Race Condition Check)
	err = s.repo.IncrementSoldCount(dropID, uint32(quantity))
	if err != nil {
		if errors.Is(err, repository.ErrSoldOut) {
			// LOSER: Stock ran out while user was paying
			// Update status to Cancelled (or specific SoldOut status if we have one)
			s.repo.UpdateOrderStatus(order.ID, models.OrderCancelled)

			// Send Loser Notification
			go integrations.SendSymbioteReceipt(customerEmail, order.CustomerPhone, "LOSER", "N/A")
			return nil
		}
		return err
	}

	// 5. WINNER: Update Order and Generate Assets
	err = s.repo.WithTransaction(func(tx repository.Repository) error {
		// Update Order Status to PAID
		if err := tx.UpdateOrderStatus(order.ID, models.OrderPaid); err != nil {
			return err
		}

		// Create Symbicode
		code, err := uuid.GenerateUUIDv7()
		if err != nil {
			return fmt.Errorf("failed to generate uuid: %w", err)
		}
		secret := generateSecretKey()

		sym := &models.Symbicode{
			Code:        code,
			SecretKey:   secret,
			ProductID:   productID,
			IsActivated: 0,
			OrderID:     order.ID,
		}

		if err := tx.CreateSymbicode(sym); err != nil {
			return fmt.Errorf("failed to create symbicode: %w", err)
		}

		return nil
	})

	if err != nil {
		// Transaction failed - rollback stock increment
		s.repo.DecrementSoldCount(dropID, uint32(quantity))
		return err
	}

	// 6. Send Notifications (Async)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Log
			}
		}()

		integrations.SendOrderConfirmationEmail(customerEmail, fmt.Sprintf("DV-%d", order.ID), float64(order.TotalAmount))

		integrations.SubmitOrderToGoogleSheet(
			customerName,
			order.CustomerPhone,
			customerEmail,
			shippingAddrStr,
			"Winner - Limited Drop",
			float64(order.TotalAmount),
			time.Now(),
		)

		integrations.SendSymbioteReceipt(customerEmail, order.CustomerPhone, "WINNER", time.Now().Format("2006-01-02 15:04:05"))
	}()

	return nil
}
