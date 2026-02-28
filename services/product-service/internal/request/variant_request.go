package request

type VariantCreateRequest struct {
	ProductID       uint    `json:"product_id" binding:"required"`
	SKU             string  `json:"sku" binding:"required"`
	BasePrice       float64 `json:"base_price" binding:"required,gt=0"`
	ComparePrice    float64 `json:"compare_price"`
	Stock           int     `json:"stock" binding:"gte=0"`
	AttributeValues []uint  `json:"attribute_values"`
}

type VariantUpdateRequest struct {
	BasePrice       float64 `json:"base_price" binding:"required,gt=0"`
	ComparePrice    float64 `json:"compare_price"`
	Stock           int     `json:"stock" binding:"gte=0"`
	AttributeValues []uint  `json:"attribute_values"`
}
