package dto

// CreateAssignmentRequest adalah payload saat trainer/admin membuat tugas.
type CreateAssignmentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Deadline    string `json:"deadline" binding:"required"` // format RFC3339, contoh: 2026-05-30T23:59:00+07:00
}

// UpdateAssignmentRequest adalah payload saat trainer/admin memperbarui tugas.
type UpdateAssignmentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Deadline    string `json:"deadline" binding:"required"` // format RFC3339
	Status      string `json:"status"`                      // optional: active, closed, archived
}
