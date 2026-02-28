package response

type AttributeValueResponse struct {
	ID          uint   `json:"id"`
	AttributeID uint   `json:"attribute_id"`
	Value       string `json:"value"`
}
