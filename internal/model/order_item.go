package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	ProductID   uuid.UUID `gorm:"type:uuid;not null" json:"product_id" binding:"required"`
	ProductCode string    `gorm:"type:varchar(20)" json:"product_code"`
	ProductName string    `gorm:"type:varchar(100)" json:"product_name"`
	QtyBox      int       `gorm:"type:integer;not null" json:"qty_box" binding:"required,min=1"`
	PricePerBox int       `gorm:"type:integer;not null" json:"price_per_box" binding:"required,min=0"`
	Subtotal    int       `gorm:"type:integer;not null" json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
	Product     *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	oi.Subtotal = oi.QtyBox * oi.PricePerBox
	return nil
}

func (oi *OrderItem) BeforeSave(tx *gorm.DB) error {
	oi.Subtotal = oi.QtyBox * oi.PricePerBox
	return nil
}
