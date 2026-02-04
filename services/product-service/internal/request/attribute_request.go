package request

type AttributeRequest struct {
	Name string `json:"name" binding:"required"`
}

type AttributePatchRequest struct {
	Name *string `json:"name,omitempty"`
}

type AttributeValueRequest struct {
	AttributeID uint   `json:"attribute_id" binding:"required"`
	Value       string `json:"value" binding:"required"`
}

type AttributeValuePatchRequest struct {
	AttributeID *uint   `json:"attribute_id,omitempty"`
	Value       *string `json:"value,omitempty"`
}
