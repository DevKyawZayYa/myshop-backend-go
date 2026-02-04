package entity

import (
	"time"

	"gorm.io/gorm"
)

type Variant struct {
	ProductID       uint             `gorm:"not null" json:"product_id"`
	Product         Product          `json:"product"`
	SKU             string           `gorm:"type:varchar(255);uniqueIndex;not null" json:"sku"`
	BasePrice       float64          `json:"base_price"`
	ComparePrice    float64          `json:"compare_price"`
	Stock           int              `json:"stock"`
	ProductImages   []ProductImage   `gorm:"foreignKey:VariantID" json:"product_images"`
	AttributeValues []AttributeValue `gorm:"many2many:variant_attribute_values;" json:"attribute_values"`

	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (v *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ComparePrice == 0 {
		v.ComparePrice = v.BasePrice
	}
	return
}
