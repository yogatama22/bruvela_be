package repository

import (
	"bruvela-backend/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatusRepository interface {
	FindAll() ([]model.PaymentStatus, error)
	FindByID(id uuid.UUID) (*model.PaymentStatus, error)
	FindByCode(code string) (*model.PaymentStatus, error)
	Create(paymentStatus *model.PaymentStatus) error
	Update(paymentStatus *model.PaymentStatus) error
	Delete(id uuid.UUID) error
}

type paymentStatusRepository struct {
	db *gorm.DB
}

func NewPaymentStatusRepository(db *gorm.DB) PaymentStatusRepository {
	return &paymentStatusRepository{db: db}
}

func (r *paymentStatusRepository) FindAll() ([]model.PaymentStatus, error) {
	var paymentStatuses []model.PaymentStatus
	err := r.db.Find(&paymentStatuses).Error
	return paymentStatuses, err
}

func (r *paymentStatusRepository) FindByID(id uuid.UUID) (*model.PaymentStatus, error) {
	var paymentStatus model.PaymentStatus
	err := r.db.Where("id = ?", id).First(&paymentStatus).Error
	if err != nil {
		return nil, err
	}
	return &paymentStatus, nil
}

func (r *paymentStatusRepository) FindByCode(code string) (*model.PaymentStatus, error) {
	var paymentStatus model.PaymentStatus
	err := r.db.Where("status_code = ?", code).First(&paymentStatus).Error
	if err != nil {
		return nil, err
	}
	return &paymentStatus, nil
}

func (r *paymentStatusRepository) Create(paymentStatus *model.PaymentStatus) error {
	return r.db.Create(paymentStatus).Error
}

func (r *paymentStatusRepository) Update(paymentStatus *model.PaymentStatus) error {
	if paymentStatus.ID == "" {
		return errors.New("ID is required")
	}
	return r.db.Save(paymentStatus).Error
}

func (r *paymentStatusRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.PaymentStatus{}, id).Error
}
