package request

type ProductImageRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	VariantID *uint  `json:"variant_id,omitempty"`
	URL       string `json:"url" binding:"required,url"`
	IsDefault bool   `json:"is_default"`
}

type ProductImagePatchRequest struct {
	VariantID *uint   `json:"variant_id,omitempty"`
	URL       *string `json:"url,omitempty" binding:"omitempty,url"`
	IsDefault *bool   `json:"is_default,omitempty"`
}
