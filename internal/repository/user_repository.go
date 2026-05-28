package repository

import (
	"context"
	"database/sql"

	"assignment-platform/internal/model"
)

// UserRepository berisi query yang berhubungan dengan tabel users.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByEmail mencari user berdasarkan email saat login.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByID mencari user berdasarkan ID dari token JWT.
func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
