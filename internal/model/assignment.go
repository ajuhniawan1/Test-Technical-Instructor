package model

import "time"

// Assignment adalah tugas yang dibuat trainer untuk suatu class.
// Satu assignment hanya dimiliki oleh satu class.
type Assignment struct {
	ID          uint64    `json:"id"`
	ClassID     uint64    `json:"class_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
	CreatedBy   uint64    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
