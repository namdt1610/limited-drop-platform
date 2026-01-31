package handlers

import (
	"ecommerce-backend/internal/integrations"
	"ecommerce-backend/internal/service"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// GetActiveDrops returns active drops
func (h *Handlers) GetActiveDrops(c fiber.Ctx) error {
	drops, err := h.service.GetActiveDrops()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch drops",
		})
	}

	return c.JSON(drops)
}

// GetDropStatus returns the status of a specific drop
func (h *Handlers) GetDropStatus(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid drop ID",
		})
	}

	status, err := h.service.GetDropStatus(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Drop not found",
		})
	}

	return c.JSON(status)
}

// PurchaseDrop handles drop purchase
func (h *Handlers) PurchaseDrop(c fiber.Ctx) error {
	idStr := c.Params("id")
	dropID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[PURCHASE] Invalid drop ID: %s\n", idStr)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid drop ID",
		})
	}

	var req struct {
		Quantity int    `json:"quantity"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Address  string `json:"address"`
		Province string `json:"province"`
		District string `json:"district"`
		Ward     string `json:"ward"`
	}

	body := c.Body()
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Fprintf(os.Stderr, "[PURCHASE] JSON unmarshal error: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body: " + err.Error(),
		})
	}

	// Validate required fields
	if err := validatePurchaseRequest(dropID, req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Create purchase request object
	purchaseReq := &service.PurchaseRequest{
		Quantity: req.Quantity,
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Address:  req.Address,
		Province: req.Province,
		District: req.District,
		Ward:     req.Ward,
	}

	result, err := h.service.PurchaseDrop(dropID, purchaseReq)
	if err != nil {
		// Do not log "Sold Out" errors as they are expected at high volume
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

// PayOSWebhook handles PayOS webhook for limited drop payments
func (h *Handlers) PayOSWebhook(c fiber.Ctx) error {
	// Get webhook signature from headers
	signature := c.Get("x-payos-signature")
	// Allow unsigned webhooks in local/dev mode when PAYOS_CLIENT_ID is not configured
	if signature == "" {
		if os.Getenv("PAYOS_CLIENT_ID") == "" {
			// dev mode: proceed without signature verification
		} else {
			return c.Status(400).JSON(fiber.Map{
				"error": "Missing webhook signature",
			})
		}
	}

	// Get raw body for signature verification
	body := c.Body()

	// Verify webhook signature when provided
	if signature != "" {
		expectedSignature := integrations.GeneratePayOSSignature(string(body))
		if signature != expectedSignature {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid webhook signature",
			})
		}
	}

	// Parse webhook payload
	var webhookData struct {
		Code string `json:"code"`
		Desc string `json:"desc"`
		Data struct {
			OrderCode     int64             `json:"orderCode"`
			Amount        int64             `json:"amount"`
			Status        string            `json:"status"`
			Description   string            `json:"description"`
			Metadata      map[string]string `json:"metadata"`
			PaymentMethod string            `json:"paymentMethod"`
		} `json:"data"`
	}

	body = c.Body()
	if err := json.Unmarshal(body, &webhookData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid webhook payload",
		})
	}

	// Only process successful payments
	if webhookData.Data.Status != "PAID" {
		return c.JSON(fiber.Map{
			"message": "Payment not completed",
		})
	}

	// Process the successful payment
	err := h.service.ProcessSuccessfulDropPayment(webhookData.Data.OrderCode)
	if err != nil {
		// Log error and return 500 to PayOS to trigger retry
		fmt.Fprintf(os.Stderr, "[WEBHOOK ERROR] Order %d: %v\n", webhookData.Data.OrderCode, err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal Server Error, please retry",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment processed successfully",
	})
}

func registerDropRoutes(app *fiber.App, h *Handlers) {
	app.Get("/api/drops", h.GetActiveDrops)
	app.Get("/api/drops/:id/status", h.GetDropStatus)
	app.Post("/api/drops/:id/purchase", h.PurchaseDrop)
	app.Post("/api/limited-drops/webhook/payos", h.PayOSWebhook)
}

func validatePurchaseRequest(dropID uint64, req struct {
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
}) error {
	if req.Name == "" {
		return fmt.Errorf("Drop %d - Họ và tên là bắt buộc", dropID)
	}
	if req.Phone == "" {
		return fmt.Errorf("Drop %d - Số điện thoại là bắt buộc", dropID)
	}
	if req.Email == "" {
		return fmt.Errorf("Drop %d - Email là bắt buộc", dropID)
	}
	if req.Address == "" {
		return fmt.Errorf("Drop %d - Địa chỉ là bắt buộc", dropID)
	}
	if req.Province == "" {
		return fmt.Errorf("Drop %d - Tỉnh / thành phố là bắt buộc", dropID)
	}
	if req.District == "" {
		return fmt.Errorf("Drop %d - Quận / huyện là bắt buộc", dropID)
	}
	if req.Ward == "" {
		return fmt.Errorf("Drop %d - Phường / xã là bắt buộc", dropID)
	}
	if req.Quantity <= 0 {
		return fmt.Errorf("Drop %d - Số lượng phải lớn hơn 0", dropID)
	}
	return nil
}
