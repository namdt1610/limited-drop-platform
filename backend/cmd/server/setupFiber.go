package main

import (
	"ecommerce-backend/config"
	"ecommerce-backend/internal/handlers"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func setupServer(cfg *config.Config, hdlrs *handlers.Handlers, db *gorm.DB) *fiber.App 