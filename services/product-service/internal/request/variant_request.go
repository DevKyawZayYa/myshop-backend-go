package request

type VariantRequest struct {
	ProductID       uint                 `json:"product_id" binding:"required"`
	SKU             string               `json:"sku" binding:"required"`
	BasePrice       float64              `json:"base_price" binding:"required,gt=0"`
	ComparePrice    float64              `json:"compare_price"`
	Stock           int                  `json:"stock" binding:"gte=0"`
	AttributeValues []uint               `json:"attribute_values"`
	ProductImages   []*ProductImageInput `json:"product_images,omitempty" binding:"omitempty,dive"`
}

type VariantPatchRequest struct {
	ProductID       *uint                 `json:"product_id,omitempty"`
	SKU             *string               `json:"sku,omitempty"`
	BasePrice       *float64              `json:"base_price,omitempty" binding:"omitempty,gt=0"`
	ComparePrice    *float64              `json:"compare_price,omitempty" binding:"omitempty,gte=0"`
	Stock           *int                  `json:"stock,omitempty" binding:"omitempty,gte=0"`
	AttributeValues *[]uint               `json:"attribute_values,omitempty"`
	ProductImages   *[]*ProductImageInput `json:"product_images,omitempty" binding:"omitempty,dive"`
}
