package request

type ProductRequest struct {
	Name          string               `json:"name" binding:"required"`
	Description   string               `json:"description,omitempty"`
	Categories    []uint               `json:"categories" binding:"required,min=1"`
	Attributes    []uint               `json:"attributes"`
	BasePrice     float64              `json:"base_price" binding:"required,gt=0"`
	ComparePrice  float64              `json:"compare_price"`
	ProductImages []*ProductImageInput `json:"product_images,omitempty" binding:"omitempty,dive"`
}

type ProductImageInput struct {
	URL       string `json:"url" binding:"required,url"`
	IsDefault bool   `json:"is_default"`
}

type ProductPatchRequest struct {
	Name          *string               `json:"name,omitempty"`
	Description   *string               `json:"description,omitempty"`
	Categories    *[]uint               `json:"categories,omitempty"`
	Attributes    *[]uint               `json:"attributes,omitempty"`
	BasePrice     *float64              `json:"base_price,omitempty" binding:"omitempty,gt=0"`
	ComparePrice  *float64              `json:"compare_price,omitempty" binding:"omitempty,gt=0"`
	ProductImages *[]*ProductImageInput `json:"product_images,omitempty" binding:"omitempty,dive"`
}
