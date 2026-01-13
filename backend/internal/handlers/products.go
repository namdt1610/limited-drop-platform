package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// ListProducts returns all products
func (h *Handlers) ListProducts(c fiber.Ctx) error {
	products, err := h.service.ListProducts()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(products)
}

// GetProduct returns a single product
func (h *Handlers) GetProduct(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.service.GetProduct(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

func registerProductRoutes(app *fiber.App, h *Handlers) {
	app.Get("/api/products", h.ListProducts)
	app.Get("/api/products/:id", h.GetProduct)
}
