package service

import (
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"
	"ecommerce-backend/internal/utils/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

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

	// 4. Execute Atomic Transaction (Stock + Order Status + Symbicode)
	err = s.repo.WithTransaction(func(tx repository.Repository) error {
		// 4.1. Increment Stock (Atomic Check)
		if err := tx.IncrementSoldCount(dropID, uint32(quantity)); err != nil {
			return err // Will be handled below (ErrSoldOut or other)
		}

		// 4.2. Update Order Status to PAID
		if err := tx.UpdateOrderStatus(order.ID, models.OrderPaid); err != nil {
			return err
		}

		// 4.3. Create Symbicode
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

	// 5. Handle Transaction Result
	if err != nil {
		if errors.Is(err, repository.ErrSoldOut) {
			// LOSER: Stock ran out during transaction attempt
			// Update status to Cancelled
			s.repo.UpdateOrderStatus(order.ID, models.OrderCancelled)

			// Send Loser Notification
			go s.email.SendSymbioteReceipt(customerEmail, order.CustomerPhone, "LOSER", "N/A")
			return nil
		}
		// Other errors: Return to retry (or log if fatal)
		return err
	}

	// 6. WINNER: Send Notifications (Async)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Log
			}
		}()

		s.email.SendOrderConfirmation(customerEmail, fmt.Sprintf("DV-%d", order.ID), float64(order.TotalAmount))

		s.sheets.SubmitOrder(
			customerName,
			order.CustomerPhone,
			customerEmail,
			shippingAddrStr,
			"Winner - Limited Drop",
			float64(order.TotalAmount),
			time.Now(),
		)

		s.email.SendSymbioteReceipt(customerEmail, order.CustomerPhone, "WINNER", time.Now().Format("2006-01-02 15:04:05"))
	}()

	return nil
}
