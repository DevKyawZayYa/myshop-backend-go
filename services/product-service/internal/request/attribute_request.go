package request

type AttributeCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type AttributeUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
