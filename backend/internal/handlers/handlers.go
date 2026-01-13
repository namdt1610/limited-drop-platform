package handlers

import (
	"ecommerce-backend/internal/service"

	"github.com/gofiber/fiber/v3"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	service service.Service
}

// NewHandlers creates new handlers instance
func NewHandlers(svc service.Service) *Handlers {
	return &Handlers{
		service: svc,
	}
}

// RegisterRoutes registers all routes by delegating to feature-specific registrars
func (h *Handlers) RegisterRoutes(app *fiber.App) {
	registerHealthRoutes(app, h)
	registerProductRoutes(app, h)
	registerDropRoutes(app, h)
	registerOrderRoutes(app, h)
	registerSymbicodeRoutes(app, h)
}
