package handlers

import "github.com/gofiber/fiber/v3"

// HealthCheck returns health status
func (h *Handlers) HealthCheck(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Service is healthy",
	})
}

func registerHealthRoutes(app *fiber.App, h *Handlers) {
	app.Get("/health", h.HealthCheck)
}
