package repository

import (
	"bruvela-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BatchRepository interface {
	Create(batch *model.Batch) error
	FindByID(id uuid.UUID) (*model.Batch, error)
	FindAll() ([]model.Batch, error)
	Update(batch *model.Batch) error
	FindOpenBatch() (*model.Batch, error)
	FindActive() (*model.Batch, error)
	FindByBatchNumber(batchNumber int) (*model.Batch, error)
	CloseAllBatches() error
	SetStatus(id uuid.UUID, status string) error
	Delete(id uuid.UUID) error
}

type batchRepository struct {
	db *gorm.DB
}

func NewBatchRepository(db *gorm.DB) BatchRepository {
	return &batchRepository{db: db}
}

func (r *batchRepository) Create(batch *model.Batch) error {
	return r.db.Create(batch).Error
}

func (r *batchRepository) FindByID(id uuid.UUID) (*model.Batch, error) {
	var batch model.Batch
	err := r.db.First(&batch, "id = ?", id).Error
	return &batch, err
}

func (r *batchRepository) FindAll() ([]model.Batch, error) {
	var batches []model.Batch
	err := r.db.Order("batch_number DESC").Find(&batches).Error
	return batches, err
}

func (r *batchRepository) Update(batch *model.Batch) error {
	return r.db.Save(batch).Error
}

func (r *batchRepository) FindOpenBatch() (*model.Batch, error) {
	var batch model.Batch
	err := r.db.Where("status = ?", "open").First(&batch).Error
	return &batch, err
}

func (r *batchRepository) FindActive() (*model.Batch, error) {
	var batch model.Batch
	err := r.db.Where("status = ?", "open").First(&batch).Error
	return &batch, err
}

func (r *batchRepository) FindByBatchNumber(batchNumber int) (*model.Batch, error) {
	var batch model.Batch
	err := r.db.Where("batch_number = ?", batchNumber).First(&batch).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *batchRepository) CloseAllBatches() error {
	return r.db.Model(&model.Batch{}).Where("status = ?", "open").Update("status", "closed").Error
}

func (r *batchRepository) SetStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Batch{}).Where("id = ?", id).Update("status", status).Error
}

func (r *batchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Batch{}, id).Error
}
