package db

import (
	"context"
	"database/sql"
)

type PostgresDB struct {
	database sql.DB
}

func (pdb *PostgresDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return pdb.database.ExecContext(ctx, query, args...)
}

func (pdb *PostgresDB) QueryRowContext(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	return pdb.database.QueryRowContext(ctx, query, args...), nil
}
func (pdb *PostgresDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return pdb.database.QueryContext(ctx, query, args...)
}
func (pdb *PostgresDB) Close() error {
	return pdb.database.Close()
}
