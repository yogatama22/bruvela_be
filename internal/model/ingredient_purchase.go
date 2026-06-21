package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Date is a custom type that handles date-only strings in JSON
type Date struct {
	time.Time
}

// UnmarshalJSON parses date-only strings (YYYY-MM-DD) or full datetime strings
func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == `""` {
		d.Time = time.Time{}
		return nil
	}

	// Remove quotes if present
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// Try parsing as date-only first
	if t, err := time.Parse("2006-01-02", s); err == nil {
		d.Time = t
		return nil
	}

	// Fall back to RFC3339
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", s)
	}
	d.Time = t
	return nil
}

// MarshalJSON outputs date in YYYY-MM-DD format
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Time.Format("2006-01-02"))), nil
}

// Value implements driver.Valuer for database
func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

// Scan implements sql.Scanner for database
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	if t, ok := value.(time.Time); ok {
		d.Time = t
		return nil
	}
	return fmt.Errorf("cannot scan %T into Date", value)
}

type IngredientPurchase struct {
	ID           uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BatchID      uuid.UUID   `gorm:"type:uuid;not null" json:"batch_id" binding:"required"`
	IngredientID uuid.UUID   `gorm:"type:uuid;not null" json:"ingredient_id" binding:"required"`
	PurchaseDate Date        `gorm:"type:date;not null" json:"purchase_date" binding:"required"`
	Supplier     string      `gorm:"type:varchar(100)" json:"supplier"`
	QtyPack      float64     `gorm:"type:decimal(10,3);not null" json:"qty_pack" binding:"required,gt=0"`
	PricePerPack int         `gorm:"type:integer;not null" json:"price_per_pack" binding:"required,min=0"`
	TotalPrice   int         `gorm:"type:integer;not null" json:"total_price"`
	Note         string      `gorm:"type:text" json:"note"`
	CreatedBy    uuid.UUID   `gorm:"type:uuid" json:"created_by"`
	CreatedAt    time.Time   `json:"created_at"`
	Batch        *Batch      `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
	Ingredient   *Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
}

func (ip *IngredientPurchase) BeforeCreate(tx *gorm.DB) error {
	if ip.ID == uuid.Nil {
		ip.ID = uuid.New()
	}
	ip.TotalPrice = int(ip.QtyPack * float64(ip.PricePerPack))
	return nil
}

func (ip *IngredientPurchase) BeforeSave(tx *gorm.DB) error {
	ip.TotalPrice = int(ip.QtyPack * float64(ip.PricePerPack))
	return nil
}
