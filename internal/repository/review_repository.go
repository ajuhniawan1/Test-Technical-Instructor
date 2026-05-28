package repository

import (
	"context"
	"database/sql"

	"assignment-platform/internal/model"
)

// ReviewRepository berisi query khusus untuk review submission.
type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// GetSubmissionForUpdate mengambil submission dan class_id dengan row-level lock.
func (r *ReviewRepository) GetSubmissionForUpdate(ctx context.Context, q DBTX, submissionID uint64) (*model.Submission, uint64, error) {
	query := `
		SELECT s.id, s.assignment_id, s.talent_id, s.github_url, s.deployment_url, s.notes, s.status,
		       s.submitted_at, s.reviewed_at, s.created_at, s.updated_at,
		       a.class_id
		FROM submissions s
		JOIN assignments a ON a.id = s.assignment_id
		WHERE s.id = ?
		FOR UPDATE
	`

	var s model.Submission
	var classID uint64
	err := q.QueryRowContext(ctx, query, submissionID).Scan(
		&s.ID,
		&s.AssignmentID,
		&s.TalentID,
		&s.GithubURL,
		&s.DeploymentURL,
		&s.Notes,
		&s.Status,
		&s.SubmittedAt,
		&s.ReviewedAt,
		&s.CreatedAt,
		&s.UpdatedAt,
		&classID,
	)
	if err != nil {
		return nil, 0, err
	}
	return &s, classID, nil
}

// InsertReview menyimpan score dan feedback trainer.
func (r *ReviewRepository) InsertReview(ctx context.Context, q DBTX, review model.SubmissionReview) error {
	query := `
		INSERT INTO submission_reviews (submission_id, trainer_id, score, feedback, review_status)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := q.ExecContext(ctx, query, review.SubmissionID, review.TrainerID, review.Score, review.Feedback, review.ReviewStatus)
	return err
}

// UpdateSubmissionStatus mengubah status submission setelah review/request revision.
func (r *ReviewRepository) UpdateSubmissionStatus(ctx context.Context, q DBTX, submissionID uint64, status string) error {
	query := `
		UPDATE submissions
		SET status = ?, reviewed_at = NOW()
		WHERE id = ?
	`
	_, err := q.ExecContext(ctx, query, status, submissionID)
	return err
}

// InsertHistory menyimpan riwayat perubahan status submission.
func (r *ReviewRepository) InsertHistory(ctx context.Context, q DBTX, history model.SubmissionHistory) error {
	query := `
		INSERT INTO submission_histories (submission_id, old_status, new_status, note, created_by)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := q.ExecContext(ctx, query, history.SubmissionID, history.OldStatus, history.NewStatus, history.Note, history.CreatedBy)
	return err
}
