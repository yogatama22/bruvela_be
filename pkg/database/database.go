package database

import (
	"bruvela-backend/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.Ingredient{},
		&model.Recipe{},
		&model.Customer{},
		&model.Batch{},
		&model.Order{},
		&model.OrderItem{},
		&model.IngredientPurchase{},
		&model.StockLog{},
		&model.JournalEntry{},
		&model.CompanySetting{},
	)
}
