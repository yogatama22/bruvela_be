package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code      string    `gorm:"type:varchar(20);unique;not null" json:"code" binding:"required"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	Price     int       `gorm:"type:integer;not null" json:"price" binding:"required,min=0"`
	PcsPerBox int       `gorm:"type:integer;default:1" json:"pcs_per_box"`
	Status    string    `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Recipes   []Recipe  `gorm:"foreignKey:ProductID" json:"recipes,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
