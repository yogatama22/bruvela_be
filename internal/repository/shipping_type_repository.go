package repository

import (
	"bruvela-backend/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShippingTypeRepository interface {
	FindAll() ([]model.ShippingType, error)
	FindByID(id uuid.UUID) (*model.ShippingType, error)
	FindByCode(code string) (*model.ShippingType, error)
	Create(shippingType *model.ShippingType) error
	Update(shippingType *model.ShippingType) error
	Delete(id uuid.UUID) error
}

type shippingTypeRepository struct {
	db *gorm.DB
}

func NewShippingTypeRepository(db *gorm.DB) ShippingTypeRepository {
	return &shippingTypeRepository{db: db}
}

func (r *shippingTypeRepository) FindAll() ([]model.ShippingType, error) {
	var shippingTypes []model.ShippingType
	err := r.db.Find(&shippingTypes).Error
	return shippingTypes, err
}

func (r *shippingTypeRepository) FindByID(id uuid.UUID) (*model.ShippingType, error) {
	var shippingType model.ShippingType
	err := r.db.Where("id = ?", id).First(&shippingType).Error
	if err != nil {
		return nil, err
	}
	return &shippingType, nil
}

func (r *shippingTypeRepository) FindByCode(code string) (*model.ShippingType, error) {
	var shippingType model.ShippingType
	err := r.db.Where("shipping_code = ?", code).First(&shippingType).Error
	if err != nil {
		return nil, err
	}
	return &shippingType, nil
}

func (r *shippingTypeRepository) Create(shippingType *model.ShippingType) error {
	return r.db.Create(shippingType).Error
}

func (r *shippingTypeRepository) Update(shippingType *model.ShippingType) error {
	if shippingType.ID == "" {
		return errors.New("ID is required")
	}
	return r.db.Save(shippingType).Error
}

func (r *shippingTypeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.ShippingType{}, id).Error
}
