package repository

import (
	"bruvela-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *model.Product) error
	FindByID(id uuid.UUID) (*model.Product, error)
	FindAll() ([]model.Product, error)
	Update(product *model.Product) error
	Delete(id uuid.UUID) error
	FindByCode(code string) (*model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uuid.UUID) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Recipes.Ingredient").First(&product, "id = ?", id).Error
	return &product, err
}

func (r *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Product{}, "id = ?", id).Error
}

func (r *productRepository) FindByCode(code string) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, "code = ?", code).Error
	return &product, err
}
