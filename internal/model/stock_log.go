package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockLog struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IngredientID  uuid.UUID   `gorm:"type:uuid;not null" json:"ingredient_id"`
	BatchID       uuid.UUID   `gorm:"type:uuid;not null" json:"batch_id"`
	LogType       string      `gorm:"type:varchar(20);not null" json:"log_type"`
	Qty           float64     `gorm:"type:decimal(10,3);not null" json:"qty"`
	StockBefore   float64     `gorm:"type:decimal(10,3)" json:"stock_before"`
	StockAfter    float64     `gorm:"type:decimal(10,3)" json:"stock_after"`
	ReferenceID   *uuid.UUID  `gorm:"type:uuid" json:"reference_id"`
	ReferenceType string      `gorm:"type:varchar(20)" json:"reference_type"`
	Note          string      `gorm:"type:text" json:"note"`
	CreatedBy     uuid.UUID   `gorm:"type:uuid" json:"created_by"`
	CreatedAt     time.Time   `json:"created_at"`
	Ingredient    *Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
	Batch         *Batch      `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
}

func (sl *StockLog) BeforeCreate(tx *gorm.DB) error {
	if sl.ID == uuid.Nil {
		sl.ID = uuid.New()
	}
	return nil
}
