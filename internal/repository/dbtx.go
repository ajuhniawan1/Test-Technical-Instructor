package repository

import (
	"context"
	"database/sql"
)

// DBTX adalah interface kecil untuk menerima *sql.DB maupun *sql.Tx.
// Dengan ini repository bisa dipakai di query biasa maupun transaction.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
