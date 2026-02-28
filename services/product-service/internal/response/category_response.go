package response

import "time"

type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ParentID    *uint     `json:"parent_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryTreeResponse struct {
	CategoryResponse
	Children []*CategoryTreeResponse `json:"children,omitempty"`
}
