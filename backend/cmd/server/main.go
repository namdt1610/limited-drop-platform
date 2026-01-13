package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce-backend/config"
	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/handlers"
	"ecommerce-backend/internal/repository"
	"ecommerce-backend/internal/service"

	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ecommerce-backend/internal/models"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize optimized SQLite database with Split Architecture (Writer/Reader)
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer database.Close()

	// Initialize GORM for migrations only (AutoMigrate)
	// Use Writer connection to avoid locking issues during migration
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect GORM for migration: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.LimitedDrop{},
		&models.Symbicode{},
	); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	log.Println("database connected, migrated and configured with Split Architecture (WAL mode)")

	// Initialize layers using optimized database (raw SQL with Writer/Reader split)
	// SmartExecutor automatically routes SELECT to Reader and writes to Writer
	executor := database.NewSmartExecutor(database.DB.Writer, database.DB.Reader)
	repo := repository.NewRepository(executor)
	svc := service.NewService(repo)
	hdlrs := handlers.NewHandlers(svc)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: gojson.Marshal,
		JSONDecoder: gojson.Unmarshal,
	})

	// Request logger
	app.Use(logger.New())

	// ETag: HTTP caching
	app.Use(etag.New())

	// Compression: Gzip/Brotli
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// CORS - Allow frontend origins
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"}
	if os.Getenv("CORS_ORIGINS") != "" {
		allowedOrigins = []string{os.Getenv("CORS_ORIGINS")}
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Register routes
	hdlrs.RegisterRoutes(app)

	// Debug route
	app.Get("/debug/products/count", func(c fiber.Ctx) error {
		var count int64
		if err := db.Model(&models.Product{}).Count(&count).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"count": count})
	})

	// Start server
	log.Printf("starting server on :%s", cfg.Port)

	// Channel to listen for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Printf("server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	log.Println("shutting down server gracefully...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("server shutdown complete")
}
