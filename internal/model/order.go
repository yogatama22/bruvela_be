package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID           uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BatchID      uuid.UUID   `gorm:"type:uuid;not null" json:"batch_id" binding:"required"`
	CustomerID   *uuid.UUID  `gorm:"type:uuid" json:"customer_id"`
	CustomerName string      `gorm:"type:varchar(100)" json:"customer_name"`
	OrderDate    time.Time   `gorm:"type:date;not null" json:"order_date" binding:"required"`
	Channel      string      `gorm:"type:varchar(30);default:'whatsapp'" json:"channel"`
	ShippingType string      `gorm:"type:varchar(20)" json:"shipping_type"`
	ShippingDest string      `gorm:"type:varchar(200)" json:"shipping_dest"`
	ShippingCost int         `gorm:"type:integer;default:0" json:"shipping_cost"`
	Discount     int         `gorm:"type:integer;default:0" json:"discount"`
	TotalProduct int         `gorm:"type:integer;default:0" json:"total_product"`
	TotalBill    int         `gorm:"type:integer;default:0" json:"total_bill"`
	PayStatus    string      `gorm:"type:varchar(20);default:'belum_bayar'" json:"pay_status"`
	ProdStatus   string      `gorm:"type:varchar(20);default:'baru'" json:"prod_status"`
	Note         string      `gorm:"type:text" json:"note"`
	CreatedBy    uuid.UUID   `gorm:"type:uuid" json:"created_by"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Batch        *Batch      `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
	Customer     *Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Items        []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (o *Order) CalculateTotals() {
	totalQty := 0
	totalPrice := 0
	for _, item := range o.Items {
		totalQty += item.QtyBox
		totalPrice += item.Subtotal
	}
	o.TotalProduct = totalQty
	o.TotalBill = totalPrice + o.ShippingCost - o.Discount
}
