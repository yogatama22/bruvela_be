package model

import "time"

type PaymentStatus struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StatusCode  string    `json:"status_code" gorm:"uniqueIndex;not null;size:20"`
	StatusName  string    `json:"status_name" gorm:"not null;size:50"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PaymentStatus) TableName() string {
	return "payment_status"
}
