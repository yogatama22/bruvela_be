package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Batch struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BatchNumber  int        `gorm:"type:integer;unique;not null" json:"batch_number" binding:"required"`
	Name         string     `gorm:"type:varchar(100)" json:"name"`
	StartDate    time.Time  `gorm:"type:date;not null" json:"start_date" binding:"required"`
	EndDate      *time.Time `gorm:"type:date" json:"end_date"`
	Status       string     `gorm:"type:varchar(20);default:'open'" json:"status"`
	TotalModal   int        `gorm:"type:integer;default:0" json:"total_modal"`
	TotalRevenue int        `gorm:"type:integer;default:0" json:"total_revenue"`
	TotalHPP     int        `gorm:"type:integer;default:0" json:"total_hpp"`
	GrossProfit  int        `gorm:"type:integer;default:0" json:"gross_profit"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (b *Batch) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (b *Batch) CalculateProfit() {
	b.GrossProfit = b.TotalRevenue - b.TotalHPP
}
