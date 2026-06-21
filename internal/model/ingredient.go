package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ingredient struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	PackUnit     string    `gorm:"type:varchar(20);not null" json:"pack_unit" binding:"required"`
	QtyPerPack   float64   `gorm:"type:decimal(10,3);not null" json:"qty_per_pack" binding:"required,gt=0"`
	UseUnit      string    `gorm:"type:varchar(20);not null" json:"use_unit" binding:"required"`
	PricePerPack int       `gorm:"type:integer;not null" json:"price_per_pack" binding:"required,min=0"`
	PricePerUse  float64   `gorm:"type:decimal(10,4)" json:"price_per_use"`
	MinStock     float64   `gorm:"type:decimal(10,3);default:0" json:"min_stock"`
	CurrentStock float64   `gorm:"type:decimal(10,3);default:0" json:"current_stock"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (i *Ingredient) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (i *Ingredient) BeforeSave(tx *gorm.DB) error {
	if i.QtyPerPack > 0 {
		i.PricePerUse = float64(i.PricePerPack) / i.QtyPerPack
	}
	return nil
}

func (i *Ingredient) IsLowStock() bool {
	return i.CurrentStock < i.MinStock
}
