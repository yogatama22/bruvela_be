package repository

import (
	"bruvela-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseRepository interface {
	Create(purchase *model.IngredientPurchase) error
	FindByID(id uuid.UUID) (*model.IngredientPurchase, error)
	FindAll(filters map[string]interface{}) ([]model.IngredientPurchase, error)
	FindByBatchID(batchID uuid.UUID) ([]model.IngredientPurchase, error)
	Update(purchase *model.IngredientPurchase) error
	Delete(id uuid.UUID) error
}

type purchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) Create(purchase *model.IngredientPurchase) error {
	return r.db.Create(purchase).Error
}

func (r *purchaseRepository) FindByID(id uuid.UUID) (*model.IngredientPurchase, error) {
	var purchase model.IngredientPurchase
	err := r.db.Preload("Ingredient").Preload("Batch").First(&purchase, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}

func (r *purchaseRepository) FindAll(filters map[string]interface{}) ([]model.IngredientPurchase, error) {
	var purchases []model.IngredientPurchase
	query := r.db.Preload("Ingredient").Preload("Batch").Order("purchase_date DESC")

	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key+" = ?", value)
		}
	}

	err := query.Find(&purchases).Error
	return purchases, err
}

func (r *purchaseRepository) FindByBatchID(batchID uuid.UUID) ([]model.IngredientPurchase, error) {
	var purchases []model.IngredientPurchase
	err := r.db.Preload("Ingredient").Preload("Batch").
		Where("batch_id = ?", batchID).
		Order("purchase_date DESC").
		Find(&purchases).Error
	return purchases, err
}

func (r *purchaseRepository) Update(purchase *model.IngredientPurchase) error {
	return r.db.Save(purchase).Error
}

func (r *purchaseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.IngredientPurchase{}, "id = ?", id).Error
}
