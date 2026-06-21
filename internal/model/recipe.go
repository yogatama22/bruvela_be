package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Recipe struct {
	ID           uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID    uuid.UUID   `gorm:"type:uuid;not null" json:"product_id" binding:"required"`
	IngredientID uuid.UUID   `gorm:"type:uuid;not null" json:"ingredient_id" binding:"required"`
	QtyPerBox    float64     `gorm:"type:decimal(10,3);not null" json:"qty_per_box" binding:"required,gt=0"`
	UseUnit      string      `gorm:"type:varchar(20);not null" json:"use_unit" binding:"required"`
	CostPerBox   float64     `gorm:"type:decimal(10,4)" json:"cost_per_box"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Product      *Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Ingredient   *Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
}

func (r *Recipe) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (r *Recipe) CalculateCost(ingredient *Ingredient) {
	if ingredient != nil {
		r.CostPerBox = r.QtyPerBox * ingredient.PricePerUse
	}
}
