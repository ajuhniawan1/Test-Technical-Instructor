package model

import "time"

// Class merepresentasikan batch/kelas bootcamp.
// Admin dapat membuat class dan assign trainer/talent ke class ini.
type Class struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
