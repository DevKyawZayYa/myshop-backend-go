package request

type ProductImageCreateRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	VariantID *uint  `json:"variant_id,omitempty"`
	URL       string `json:"url" binding:"required,url"`
	IsDefault bool   `json:"is_default"`
}

type ProductImageUpdateRequest struct {
	URL       string `json:"url" binding:"required,url"`
	IsDefault bool   `json:"is_default"`
}
