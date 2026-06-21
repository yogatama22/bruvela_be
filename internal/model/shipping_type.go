package model

import "time"

type ShippingType struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ShippingCode string    `json:"shipping_code" gorm:"uniqueIndex;not null;size:20"`
	ShippingName string    `json:"shipping_name" gorm:"not null;size:100"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ShippingType) TableName() string {
	return "shipping_type"
}
