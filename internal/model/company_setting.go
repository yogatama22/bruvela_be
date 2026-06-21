package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanySetting struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyName  string    `gorm:"type:varchar(200);not null" json:"company_name"`
	Address      string    `gorm:"type:text" json:"address"`
	PhoneNumber  string    `gorm:"type:varchar(50)" json:"phone_number"`
	Email        string    `gorm:"type:varchar(100)" json:"email"`
	Instagram    string    `gorm:"type:varchar(100)" json:"instagram"`
	WhatsApp     string    `gorm:"type:varchar(50)" json:"whatsapp"`
	TikTok       string    `gorm:"type:varchar(100)" json:"tiktok"`
	LogoURL      string    `gorm:"type:varchar(500)" json:"logo_url"`
	InvoiceNote  string    `gorm:"type:text" json:"invoice_note"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (CompanySetting) TableName() string {
	return "company_settings"
}

func (cs *CompanySetting) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	return nil
}
