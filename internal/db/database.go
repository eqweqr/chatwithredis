package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type UserDB interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) (*sql.Row, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	Close() error
}
