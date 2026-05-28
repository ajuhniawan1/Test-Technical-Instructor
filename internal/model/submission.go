package model

import (
	"database/sql"
	"time"
)

// Submission adalah data pengumpulan tugas dari talent.
// Status digunakan untuk melacak proses: submitted, reviewed, revision_required, late.
type Submission struct {
	ID            uint64       `json:"id"`
	AssignmentID  uint64       `json:"assignment_id"`
	TalentID      uint64       `json:"talent_id"`
	GithubURL     string       `json:"github_url"`
	DeploymentURL string       `json:"deployment_url"`
	Notes         string       `json:"notes"`
	Status        string       `json:"status"`
	SubmittedAt   sql.NullTime `json:"-"`
	ReviewedAt    sql.NullTime `json:"-"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// SubmissionReview menyimpan score dan feedback trainer.
type SubmissionReview struct {
	ID           uint64    `json:"id"`
	SubmissionID uint64    `json:"submission_id"`
	TrainerID    uint64    `json:"trainer_id"`
	Score        int       `json:"score"`
	Feedback     string    `json:"feedback"`
	ReviewStatus string    `json:"review_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SubmissionHistory menyimpan riwayat perubahan status submission.
// Ini membantu tracking saat ada revisi atau review.
type SubmissionHistory struct {
	ID           uint64    `json:"id"`
	SubmissionID uint64    `json:"submission_id"`
	OldStatus    string    `json:"old_status"`
	NewStatus    string    `json:"new_status"`
	Note         string    `json:"note"`
	CreatedBy    uint64    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}
