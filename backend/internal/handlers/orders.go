package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// GetOrderByID retrieves a specific order by ID
func (h *Handlers) GetOrderByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := h.service.GetOrderByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return c.JSON(order)
}

// GetOrdersByPhone retrieves all orders for a user (by phone from query param)
func (h *Handlers) GetOrdersByPhone(c fiber.Ctx) error {
	phone := c.Query("phone")
	if phone == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Phone parameter is required",
		})
	}

	orders, err := h.service.GetOrdersByUserPhone(phone)
	if err != nil {
		// Log the error
		println("Error retrieving orders: " + err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve orders: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"count":  len(orders),
	})
}

func registerOrderRoutes(app *fiber.App, h *Handlers) {
	app.Get("/api/orders/:id", h.GetOrderByID)
	app.Get("/api/orders", h.GetOrdersByPhone)
}
