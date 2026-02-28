package response

import "time"

type AttributeResponse struct {
	ID        uint                     `json:"id"`
	Name      string                   `json:"name"`
	Values    []AttributeValueResponse `json:"values,omitempty"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}
