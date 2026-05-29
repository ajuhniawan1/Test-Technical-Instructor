package repository

import (
	"context"
	"database/sql"
	"time"

	"assignment-platform/internal/model"
)

// ClassRepository berisi query untuk class, class_trainers, dan class_talents.
type ClassRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

// Create membuat class/batch baru.
func (r *ClassRepository) Create(ctx context.Context, class model.Class) (uint64, error) {
	query := `
		INSERT INTO classes (name, description, start_date, end_date)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query, class.Name, class.Description, class.StartDate, class.EndDate)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return uint64(id), err
}

// List mengambil semua class tanpa pagination. Fungsi ini tetap disediakan untuk kebutuhan internal/test.
func (r *ClassRepository) List(ctx context.Context) ([]model.Class, error) {
	query := `
		SELECT id, name, description, start_date, end_date, created_at, updated_at
		FROM classes
		ORDER BY id DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classes := make([]model.Class, 0)
	for rows.Next() {
		var class model.Class
		if err := rows.Scan(&class.ID, &class.Name, &class.Description, &class.StartDate, &class.EndDate, &class.CreatedAt, &class.UpdatedAt); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, rows.Err()
}

// ListPaginated mengambil daftar class dengan limit dan offset.
func (r *ClassRepository) ListPaginated(ctx context.Context, limit int, offset int) ([]model.Class, int, error) {
	countQuery := `SELECT COUNT(1) FROM classes`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, description, start_date, end_date, created_at, updated_at
		FROM classes
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	classes := make([]model.Class, 0)
	for rows.Next() {
		var class model.Class
		if err := rows.Scan(&class.ID, &class.Name, &class.Description, &class.StartDate, &class.EndDate, &class.CreatedAt, &class.UpdatedAt); err != nil {
			return nil, 0, err
		}
		classes = append(classes, class)
	}

	return classes, total, rows.Err()
}

// FindByID mengambil detail class.
func (r *ClassRepository) FindByID(ctx context.Context, classID uint64) (*model.Class, error) {
	query := `
		SELECT id, name, description, start_date, end_date, created_at, updated_at
		FROM classes
		WHERE id = ?
		LIMIT 1
	`
	var class model.Class
	err := r.db.QueryRowContext(ctx, query, classID).Scan(
		&class.ID,
		&class.Name,
		&class.Description,
		&class.StartDate,
		&class.EndDate,
		&class.CreatedAt,
		&class.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &class, nil
}

// AssignTrainer menghubungkan trainer ke class.
func (r *ClassRepository) AssignTrainer(ctx context.Context, classID uint64, trainerID uint64) error {
	query := `
		INSERT INTO class_trainers (class_id, trainer_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE trainer_id = trainer_id
	`
	_, err := r.db.ExecContext(ctx, query, classID, trainerID)
	return err
}

// AssignTalent menghubungkan talent ke class.
func (r *ClassRepository) AssignTalent(ctx context.Context, classID uint64, talentID uint64) error {
	query := `
		INSERT INTO class_talents (class_id, talent_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE talent_id = talent_id
	`
	_, err := r.db.ExecContext(ctx, query, classID, talentID)
	return err
}

// IsTrainerAssigned mengecek apakah trainer handle class tertentu.
func (r *ClassRepository) IsTrainerAssigned(ctx context.Context, classID uint64, trainerID uint64) (bool, error) {
	query := `SELECT COUNT(1) FROM class_trainers WHERE class_id = ? AND trainer_id = ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, classID, trainerID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsTalentAssigned mengecek apakah talent terdaftar di class tertentu.
func (r *ClassRepository) IsTalentAssigned(ctx context.Context, classID uint64, talentID uint64) (bool, error) {
	query := `SELECT COUNT(1) FROM class_talents WHERE class_id = ? AND talent_id = ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, classID, talentID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountTalentsInClass menghitung jumlah talent dalam class.
func (r *ClassRepository) CountTalentsInClass(ctx context.Context, classID uint64) (int, error) {
	query := `SELECT COUNT(1) FROM class_talents WHERE class_id = ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, classID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// CountClassesForTalent menghitung jumlah class yang diikuti talent.
func (r *ClassRepository) CountClassesForTalent(ctx context.Context, talentID uint64) (int, error) {
	query := `SELECT COUNT(1) FROM class_talents WHERE talent_id = ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, talentID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// ParseDate adalah helper sederhana untuk parsing tanggal class.
func ParseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

// ListTrainersByClass mengambil daftar trainer yang terdaftar pada class tertentu.
func (r *ClassRepository) ListTrainersByClass(ctx context.Context, classID uint64) ([]model.ClassTrainerMember, error) {
	query := `
		SELECT
			ct.id,
			ct.class_id,
			ct.trainer_id,
			u.name AS trainer_name,
			u.email AS trainer_email,
			u.role AS trainer_role,
			ct.created_at
		FROM class_trainers ct
		JOIN users u ON u.id = ct.trainer_id
		WHERE ct.class_id = ?
		ORDER BY ct.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.ClassTrainerMember, 0)

	for rows.Next() {
		var item model.ClassTrainerMember

		if err := rows.Scan(
			&item.ID,
			&item.ClassID,
			&item.TrainerID,
			&item.TrainerName,
			&item.TrainerEmail,
			&item.TrainerRole,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

// ListTalentsByClass mengambil daftar talent yang terdaftar pada class tertentu.
func (r *ClassRepository) ListTalentsByClass(ctx context.Context, classID uint64) ([]model.ClassTalentMember, error) {
	query := `
		SELECT
			ct.id,
			ct.class_id,
			ct.talent_id,
			u.name AS talent_name,
			u.email AS talent_email,
			u.role AS talent_role,
			ct.created_at
		FROM class_talents ct
		JOIN users u ON u.id = ct.talent_id
		WHERE ct.class_id = ?
		ORDER BY ct.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.ClassTalentMember, 0)

	for rows.Next() {
		var item model.ClassTalentMember

		if err := rows.Scan(
			&item.ID,
			&item.ClassID,
			&item.TalentID,
			&item.TalentName,
			&item.TalentEmail,
			&item.TalentRole,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}
