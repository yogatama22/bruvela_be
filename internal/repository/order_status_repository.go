package repository

import (
	"bruvela-backend/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatusRepository interface {
	FindAll() ([]model.OrderStatus, error)
	FindByID(id uuid.UUID) (*model.OrderStatus, error)
	FindByCode(code string) (*model.OrderStatus, error)
	Create(orderStatus *model.OrderStatus) error
	Update(orderStatus *model.OrderStatus) error
	Delete(id uuid.UUID) error
}

type orderStatusRepository struct {
	db *gorm.DB
}

func NewOrderStatusRepository(db *gorm.DB) OrderStatusRepository {
	return &orderStatusRepository{db: db}
}

func (r *orderStatusRepository) FindAll() ([]model.OrderStatus, error) {
	var orderStatuses []model.OrderStatus
	err := r.db.Find(&orderStatuses).Error
	return orderStatuses, err
}

func (r *orderStatusRepository) FindByID(id uuid.UUID) (*model.OrderStatus, error) {
	var orderStatus model.OrderStatus
	err := r.db.Where("id = ?", id).First(&orderStatus).Error
	if err != nil {
		return nil, err
	}
	return &orderStatus, nil
}

func (r *orderStatusRepository) FindByCode(code string) (*model.OrderStatus, error) {
	var orderStatus model.OrderStatus
	err := r.db.Where("order_status_code = ?", code).First(&orderStatus).Error
	if err != nil {
		return nil, err
	}
	return &orderStatus, nil
}

func (r *orderStatusRepository) Create(orderStatus *model.OrderStatus) error {
	return r.db.Create(orderStatus).Error
}

func (r *orderStatusRepository) Update(orderStatus *model.OrderStatus) error {
	if orderStatus.ID == "" {
		return errors.New("ID is required")
	}
	return r.db.Save(orderStatus).Error
}

func (r *orderStatusRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.OrderStatus{}, id).Error
}
