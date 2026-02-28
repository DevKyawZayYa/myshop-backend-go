package response

import "time"

type ProductResponse struct {
	ID           uint                 `json:"id"`
	Name         string               `json:"name"`
	Description  string               `json:"description,omitempty"`
	BasePrice    float64              `json:"base_price"`
	ComparePrice float64              `json:"compare_price"`
	Categories   []CategoryResponse   `json:"categories,omitempty"`
	Attributes   []AttributeResponse  `json:"attributes,omitempty"`
	Variants     []VariantResponse    `json:"variants,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}
