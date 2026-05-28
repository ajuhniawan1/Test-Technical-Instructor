package model

import "time"

// User merepresentasikan data pengguna sistem.
// Role menentukan hak akses user: admin, trainer, atau talent.
type User struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // tanda '-' agar password_hash tidak muncul di response JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
