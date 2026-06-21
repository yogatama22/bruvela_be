package repository

import (
	"bruvela-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	Create(ingredient *model.Ingredient) error
	FindByID(id uuid.UUID) (*model.Ingredient, error)
	FindAll() ([]model.Ingredient, error)
	Update(ingredient *model.Ingredient) error
	Delete(id uuid.UUID) error
	FindLowStock() ([]model.Ingredient, error)
	UpdateStock(id uuid.UUID, qty float64) error
}

type ingredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) IngredientRepository {
	return &ingredientRepository{db: db}
}

func (r *ingredientRepository) Create(ingredient *model.Ingredient) error {
	return r.db.Create(ingredient).Error
}

func (r *ingredientRepository) FindByID(id uuid.UUID) (*model.Ingredient, error) {
	var ingredient model.Ingredient
	err := r.db.First(&ingredient, "id = ?", id).Error
	return &ingredient, err
}

func (r *ingredientRepository) FindAll() ([]model.Ingredient, error) {
	var ingredients []model.Ingredient
	err := r.db.Find(&ingredients).Error
	return ingredients, err
}

func (r *ingredientRepository) Update(ingredient *model.Ingredient) error {
	return r.db.Save(ingredient).Error
}

func (r *ingredientRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Ingredient{}, "id = ?", id).Error
}

func (r *ingredientRepository) FindLowStock() ([]model.Ingredient, error) {
	var ingredients []model.Ingredient
	err := r.db.Where("current_stock < min_stock").Find(&ingredients).Error
	return ingredients, err
}

func (r *ingredientRepository) UpdateStock(id uuid.UUID, qty float64) error {
	return r.db.Model(&model.Ingredient{}).Where("id = ?", id).
		Update("current_stock", gorm.Expr("current_stock + ?", qty)).Error
}
