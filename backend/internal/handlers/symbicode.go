package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

// VerifySymbicode verifies a symbicode
func (h *Handlers) VerifySymbicode(c fiber.Ctx) error {
	var req struct {
		Code string `json:"code"`
	}

	body := c.Body()
	if err := json.Unmarshal(body, &req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	symbicode, isFirst, err := h.service.VerifySymbicode(req.Code)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid symbicode",
		})
	}

	return c.JSON(fiber.Map{
		"symbicode":           symbicode,
		"is_first_activation": isFirst,
	})
}

func registerSymbicodeRoutes(app *fiber.App, h *Handlers) {
	app.Post("/api/symbicode/verify", h.VerifySymbicode)
}
