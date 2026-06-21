package repository

import (
	"bruvela-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockLogRepository interface {
	Create(tx *gorm.DB, log *model.StockLog) error
	FindAll(filters map[string]interface{}) ([]model.StockLog, error)
	FindByIngredient(ingredientID uuid.UUID) ([]model.StockLog, error)
	FindByBatch(batchID uuid.UUID) ([]model.StockLog, error)
}

type stockLogRepository struct {
	db *gorm.DB
}

func NewStockLogRepository(db *gorm.DB) StockLogRepository {
	return &stockLogRepository{db: db}
}

func (r *stockLogRepository) Create(tx *gorm.DB, log *model.StockLog) error {
	return tx.Create(log).Error
}

func (r *stockLogRepository) FindAll(filters map[string]interface{}) ([]model.StockLog, error) {
	var logs []model.StockLog
	query := r.db.Preload("Ingredient").Preload("Batch").Order("created_at DESC")

	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key+" = ?", value)
		}
	}

	err := query.Find(&logs).Error
	return logs, err
}

func (r *stockLogRepository) FindByIngredient(ingredientID uuid.UUID) ([]model.StockLog, error) {
	var logs []model.StockLog
	err := r.db.Preload("Ingredient").Preload("Batch").
		Where("ingredient_id = ?", ingredientID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *stockLogRepository) FindByBatch(batchID uuid.UUID) ([]model.StockLog, error) {
	var logs []model.StockLog
	err := r.db.Preload("Ingredient").Preload("Batch").
		Where("batch_id = ?", batchID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}
