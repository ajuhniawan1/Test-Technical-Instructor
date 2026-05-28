package repository

import (
	"context"
	"database/sql"

	"assignment-platform/internal/model"
)

// SubmissionRepository berisi query submission dan progress.
type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

// FindByAssignmentAndTalent mencari submission milik talent pada assignment tertentu.
// Fungsi ini menerima DBTX agar bisa dipakai di luar atau di dalam transaction.
func (r *SubmissionRepository) FindByAssignmentAndTalent(ctx context.Context, q DBTX, assignmentID uint64, talentID uint64) (*model.Submission, error) {
	query := `
		SELECT id, assignment_id, talent_id, github_url, deployment_url, notes, status, submitted_at, reviewed_at, created_at, updated_at
		FROM submissions
		WHERE assignment_id = ? AND talent_id = ?
		LIMIT 1
	`

	var s model.Submission
	err := q.QueryRowContext(ctx, query, assignmentID, talentID).Scan(
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
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSubmission membuat submission baru.
func (r *SubmissionRepository) CreateSubmission(ctx context.Context, q DBTX, submission model.Submission) (uint64, error) {
	query := `
		INSERT INTO submissions (assignment_id, talent_id, github_url, deployment_url, notes, status, submitted_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`
	result, err := q.ExecContext(ctx, query,
		submission.AssignmentID,
		submission.TalentID,
		submission.GithubURL,
		submission.DeploymentURL,
		submission.Notes,
		submission.Status,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}

// UpdateSubmissionLink dipakai saat talent resubmit setelah revision_required.
func (r *SubmissionRepository) UpdateSubmissionLink(ctx context.Context, q DBTX, submissionID uint64, githubURL string, deploymentURL string, notes string, status string) error {
	query := `
		UPDATE submissions
		SET github_url = ?, deployment_url = ?, notes = ?, status = ?, submitted_at = NOW(), reviewed_at = NULL
		WHERE id = ?
	`
	_, err := q.ExecContext(ctx, query, githubURL, deploymentURL, notes, status, submissionID)
	return err
}

// GetSubmissionForUpdate mengambil submission dengan row-level lock.
// SELECT ... FOR UPDATE mencegah dua proses mengubah row yang sama secara bersamaan.
func (r *SubmissionRepository) GetSubmissionForUpdate(ctx context.Context, q DBTX, submissionID uint64) (*model.Submission, uint64, error) {
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

// InsertReview tetap disediakan untuk kompatibilitas jika review digabung dengan submission service.
func (r *SubmissionRepository) InsertReview(ctx context.Context, q DBTX, review model.SubmissionReview) error {
	query := `
		INSERT INTO submission_reviews (submission_id, trainer_id, score, feedback, review_status)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := q.ExecContext(ctx, query, review.SubmissionID, review.TrainerID, review.Score, review.Feedback, review.ReviewStatus)
	return err
}

// UpdateStatusAfterReview mengubah status submission setelah direview.
func (r *SubmissionRepository) UpdateStatusAfterReview(ctx context.Context, q DBTX, submissionID uint64, status string) error {
	query := `
		UPDATE submissions
		SET status = ?, reviewed_at = NOW()
		WHERE id = ?
	`
	_, err := q.ExecContext(ctx, query, status, submissionID)
	return err
}

// InsertHistory menyimpan riwayat perubahan status.
func (r *SubmissionRepository) InsertHistory(ctx context.Context, q DBTX, history model.SubmissionHistory) error {
	query := `
		INSERT INTO submission_histories (submission_id, old_status, new_status, note, created_by)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := q.ExecContext(ctx, query, history.SubmissionID, history.OldStatus, history.NewStatus, history.Note, history.CreatedBy)
	return err
}

// ListByTalent mengambil semua submission milik talent yang login.
func (r *SubmissionRepository) ListByTalent(ctx context.Context, talentID uint64) ([]map[string]any, error) {
	query := `
		SELECT s.id, s.assignment_id, a.title, s.github_url, s.deployment_url, s.notes, s.status, s.submitted_at, s.reviewed_at
		FROM submissions s
		JOIN assignments a ON a.id = s.assignment_id
		WHERE s.talent_id = ?
		ORDER BY s.updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, talentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, assignmentID uint64
		var title, githubURL, deploymentURL, notes, status string
		var submittedAt, reviewedAt sql.NullTime
		if err := rows.Scan(&id, &assignmentID, &title, &githubURL, &deploymentURL, &notes, &status, &submittedAt, &reviewedAt); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{
			"id":               id,
			"assignment_id":    assignmentID,
			"assignment_title": title,
			"github_url":       githubURL,
			"deployment_url":   deploymentURL,
			"notes":            notes,
			"status":           status,
			"submitted_at":     submittedAt,
			"reviewed_at":      reviewedAt,
		})
	}
	return items, rows.Err()
}

// ListByClass mengambil submission berdasarkan class untuk trainer/admin.
func (r *SubmissionRepository) ListByClass(ctx context.Context, classID uint64) ([]map[string]any, error) {
	query := `
		SELECT s.id, s.assignment_id, u.id AS talent_id, u.name AS talent_name, a.title AS assignment_title,
		       s.github_url, s.deployment_url, s.notes, s.status, s.submitted_at, s.reviewed_at
		FROM submissions s
		JOIN users u ON u.id = s.talent_id
		JOIN assignments a ON a.id = s.assignment_id
		WHERE a.class_id = ?
		ORDER BY s.updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, assignmentID, talentID uint64
		var talentName, assignmentTitle, githubURL, deploymentURL, notes, status string
		var submittedAt, reviewedAt sql.NullTime
		if err := rows.Scan(&id, &assignmentID, &talentID, &talentName, &assignmentTitle, &githubURL, &deploymentURL, &notes, &status, &submittedAt, &reviewedAt); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{
			"id":               id,
			"assignment_id":    assignmentID,
			"talent_id":        talentID,
			"talent_name":      talentName,
			"assignment_title": assignmentTitle,
			"github_url":       githubURL,
			"deployment_url":   deploymentURL,
			"notes":            notes,
			"status":           status,
			"submitted_at":     submittedAt,
			"reviewed_at":      reviewedAt,
		})
	}
	return items, rows.Err()
}

// CountStatusByClass menghitung jumlah submission per status di satu class.
func (r *SubmissionRepository) CountStatusByClass(ctx context.Context, classID uint64) (map[string]int, int, error) {
	query := `
		SELECT s.status, COUNT(1)
		FROM submissions s
		JOIN assignments a ON a.id = s.assignment_id
		WHERE a.class_id = ?
		GROUP BY s.status
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result := map[string]int{
		"submitted":         0,
		"reviewed":          0,
		"revision_required": 0,
		"late":              0,
	}

	totalSubmittedRecords := 0
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, 0, err
		}
		result[status] = count
		totalSubmittedRecords += count
	}

	return result, totalSubmittedRecords, rows.Err()
}

// CountStatusByTalent menghitung jumlah submission per status untuk satu talent.
func (r *SubmissionRepository) CountStatusByTalent(ctx context.Context, talentID uint64) (map[string]int, int, error) {
	query := `
		SELECT status, COUNT(1)
		FROM submissions
		WHERE talent_id = ?
		GROUP BY status
	`
	rows, err := r.db.QueryContext(ctx, query, talentID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result := map[string]int{
		"submitted":         0,
		"reviewed":          0,
		"revision_required": 0,
		"late":              0,
	}

	totalSubmittedRecords := 0
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, 0, err
		}
		result[status] = count
		totalSubmittedRecords += count
	}

	return result, totalSubmittedRecords, rows.Err()
}
