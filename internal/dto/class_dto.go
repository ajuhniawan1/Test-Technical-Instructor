package dto

// CreateClassRequest adalah request untuk membuat class/batch baru.
type CreateClassRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	StartDate   string `json:"start_date" binding:"required"` // format: YYYY-MM-DD
	EndDate     string `json:"end_date" binding:"required"`   // format: YYYY-MM-DD
}

// AssignUserRequest dipakai untuk assign trainer/talent ke class.
type AssignUserRequest struct {
	UserID uint64 `json:"user_id" binding:"required"`
}

// PaginationMeta menjelaskan metadata pagination pada list endpoint.
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse adalah response list data dengan pagination.
type PaginatedResponse struct {
	Items any            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}
