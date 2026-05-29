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

// ClassTrainerMember merepresentasikan data trainer yang terdaftar di sebuah class.
type ClassTrainerMember struct {
	ID           uint64    `json:"id"`
	ClassID      uint64    `json:"class_id"`
	TrainerID    uint64    `json:"trainer_id"`
	TrainerName  string    `json:"trainer_name"`
	TrainerEmail string    `json:"trainer_email"`
	TrainerRole  string    `json:"trainer_role"`
	CreatedAt    time.Time `json:"created_at"`
}

// ClassTalentMember merepresentasikan data talent yang terdaftar di sebuah class.
type ClassTalentMember struct {
	ID          uint64    `json:"id"`
	ClassID     uint64    `json:"class_id"`
	TalentID    uint64    `json:"talent_id"`
	TalentName  string    `json:"talent_name"`
	TalentEmail string    `json:"talent_email"`
	TalentRole  string    `json:"talent_role"`
	CreatedAt   time.Time `json:"created_at"`
}
