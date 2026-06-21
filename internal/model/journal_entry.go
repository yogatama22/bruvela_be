package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JournalEntry struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BatchID     uuid.UUID  `gorm:"type:uuid;not null" json:"batch_id" binding:"required"`
	EntryDate   time.Time  `gorm:"type:date;not null" json:"entry_date" binding:"required"`
	Description string     `gorm:"type:varchar(255);not null" json:"description" binding:"required"`
	Type        string     `gorm:"type:varchar(20);not null" json:"type" binding:"required"`
	Amount      int        `gorm:"type:integer;not null" json:"amount" binding:"required"`
	Balance     int        `gorm:"type:integer;default:0" json:"balance"`
	Partner     string     `gorm:"type:varchar(100)" json:"partner"`
	ReferenceID *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	Batch       *Batch     `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
}

func (je *JournalEntry) BeforeCreate(tx *gorm.DB) error {
	if je.ID == uuid.Nil {
		je.ID = uuid.New()
	}
	return nil
}
