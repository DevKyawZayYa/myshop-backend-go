package params

type BaseQueryParam struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type ProductQueryParam struct {
	BaseQueryParam
	CategoryID string `form:"category_id"`
}

func NewProductQueryParam() *ProductQueryParam {
	return &ProductQueryParam{
		BaseQueryParam: BaseQueryParam{
			Limit:  10,
			Offset: 0,
		},
	}
}

type VariantQueryParam struct {
	BaseQueryParam
	ProductID string `form:"product_id"`
}

func NewVariantQueryParam() *VariantQueryParam {
	return &VariantQueryParam{
		BaseQueryParam: BaseQueryParam{
			Limit:  10,
			Offset: 0,
		},
	}
}
