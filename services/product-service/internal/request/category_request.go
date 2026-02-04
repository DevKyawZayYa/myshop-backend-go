package request

type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
}

type CategoryPatchRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ParentID    *uint   `json:"parent_id,omitempty"`
}
