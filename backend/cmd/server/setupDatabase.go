package main

import (
	"ecommerce-backend/config"
	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func bootstrapDatabase(cfg *config.Config) (*gorm.DB, error) {
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.LimitedDrop{},
		&models.Symbicode{},
	); err != nil {
		return nil, err
	}

	log.Println("database connected & migrated successfully")
	return db, nil
}
