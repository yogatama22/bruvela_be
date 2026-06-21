package model

import "time"

type OrderStatus struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderStatusCode string   `json:"order_status_code" gorm:"uniqueIndex;not null;size:20"`
	OrderStatusName string   `json:"order_status_name" gorm:"not null;size:50"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (OrderStatus) TableName() string {
	return "order_status"
}
