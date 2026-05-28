package repository

import (
	"context"
	"database/sql"

	"assignment-platform/internal/model"
)

// AssignmentRepository berisi query untuk tabel assignments.
type AssignmentRepository struct {
	db *sql.DB
}

func NewAssignmentRepository(db *sql.DB) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

// Create menyimpan assignment baru ke database.
func (r *AssignmentRepository) Create(ctx context.Context, assignment model.Assignment) (uint64, error) {
	query := `
		INSERT INTO assignments (class_id, title, description, deadline, status, created_by)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		assignment.ClassID,
		assignment.Title,
		assignment.Description,
		assignment.Deadline,
		assignment.Status,
		assignment.CreatedBy,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}

// ListByClass mengambil assignment berdasarkan class_id.
func (r *AssignmentRepository) ListByClass(ctx context.Context, classID uint64) ([]model.Assignment, error) {
	query := `
		SELECT id, class_id, title, description, deadline, status, created_by, created_at, updated_at
		FROM assignments
		WHERE class_id = ?
		ORDER BY deadline ASC
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := make([]model.Assignment, 0)
	for rows.Next() {
		var a model.Assignment
		if err := rows.Scan(&a.ID, &a.ClassID, &a.Title, &a.Description, &a.Deadline, &a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}

// FindByID mengambil detail assignment.
func (r *AssignmentRepository) FindByID(ctx context.Context, assignmentID uint64) (*model.Assignment, error) {
	query := `
		SELECT id, class_id, title, description, deadline, status, created_by, created_at, updated_at
		FROM assignments
		WHERE id = ?
		LIMIT 1
	`
	var a model.Assignment
	err := r.db.QueryRowContext(ctx, query, assignmentID).Scan(
		&a.ID,
		&a.ClassID,
		&a.Title,
		&a.Description,
		&a.Deadline,
		&a.Status,
		&a.CreatedBy,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Update mengubah data assignment.
func (r *AssignmentRepository) Update(ctx context.Context, assignment model.Assignment) error {
	query := `
		UPDATE assignments
		SET title = ?, description = ?, deadline = ?, status = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		assignment.Title,
		assignment.Description,
		assignment.Deadline,
		assignment.Status,
		assignment.ID,
	)
	return err
}

// Close mengubah status assignment menjadi closed.
func (r *AssignmentRepository) Close(ctx context.Context, assignmentID uint64) error {
	query := `UPDATE assignments SET status = 'closed' WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, assignmentID)
	return err
}

// CountByClass menghitung jumlah assignment di satu class.
func (r *AssignmentRepository) CountByClass(ctx context.Context, classID uint64) (int, error) {
	query := `SELECT COUNT(1) FROM assignments WHERE class_id = ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, classID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// CountForTalent menghitung total assignment dari semua class yang diikuti talent.
func (r *AssignmentRepository) CountForTalent(ctx context.Context, talentID uint64) (int, error) {
	query := `
		SELECT COUNT(1)
		FROM assignments a
		JOIN class_talents ct ON ct.class_id = a.class_id
		WHERE ct.talent_id = ?
	`
	var count int
	if err := r.db.QueryRowContext(ctx, query, talentID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
