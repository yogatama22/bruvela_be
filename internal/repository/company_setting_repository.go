package repository

import (
	"bruvela-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanySettingRepository interface {
	Get() (*model.CompanySetting, error)
	Upsert(setting *model.CompanySetting) error
}

type companySettingRepository struct {
	db *gorm.DB
}

func NewCompanySettingRepository(db *gorm.DB) CompanySettingRepository {
	return &companySettingRepository{db: db}
}

func (r *companySettingRepository) Get() (*model.CompanySetting, error) {
	var setting model.CompanySetting
	err := r.db.First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *companySettingRepository) Upsert(setting *model.CompanySetting) error {
	var existing model.CompanySetting
	err := r.db.First(&existing).Error
	if err != nil {
		setting.ID = uuid.New()
		return r.db.Create(setting).Error
	}
	setting.ID = existing.ID
	return r.db.Save(setting).Error
}
