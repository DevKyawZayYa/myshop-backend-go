package request

type ProductCreateRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	Categories   []uint  `json:"categories" binding:"required,min=1"`
	Attributes   []uint  `json:"attributes"`
	BasePrice    float64 `json:"base_price" binding:"required,gt=0"`
	ComparePrice float64 `json:"compare_price"`
}

type ProductUpdateRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	Categories   []uint  `json:"categories" binding:"required,min=1"`
	Attributes   []uint  `json:"attributes"`
	BasePrice    float64 `json:"base_price" binding:"required,gt=0"`
	ComparePrice float64 `json:"compare_price"`
}
