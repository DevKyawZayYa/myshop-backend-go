package request

type AttributeValueCreateRequest struct {
    AttributeID uint   `json:"attribute_id"`
    Value       string `json:"value"`
}

type AttributeValueUpdateRequest struct {
    Value string `json:"value"`
}
