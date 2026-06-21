package repository

import (
	"bruvela-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	Create(recipe *model.Recipe) error
	FindByProductID(productID uuid.UUID) ([]model.Recipe, error)
	Update(recipe *model.Recipe) error
	Delete(id uuid.UUID) error
	DeleteByProductID(productID uuid.UUID) error
}

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) Create(recipe *model.Recipe) error {
	return r.db.Create(recipe).Error
}

func (r *recipeRepository) FindByProductID(productID uuid.UUID) ([]model.Recipe, error) {
	var recipes []model.Recipe
	err := r.db.Preload("Ingredient").Where("product_id = ?", productID).Find(&recipes).Error
	return recipes, err
}

func (r *recipeRepository) Update(recipe *model.Recipe) error {
	return r.db.Save(recipe).Error
}

func (r *recipeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Recipe{}, "id = ?", id).Error
}

func (r *recipeRepository) DeleteByProductID(productID uuid.UUID) error {
	return r.db.Where("product_id = ?", productID).Delete(&model.Recipe{}).Error
}
