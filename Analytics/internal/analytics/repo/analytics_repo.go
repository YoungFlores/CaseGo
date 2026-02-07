package analyticsRepo

import (
	"context"
	"database/sql"
)

type AnalyticsRepo interface {
	
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
    PrepareContext(context.Context, string) (*sql.Stmt, error)
    QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type PostgresAnalyticRepo struct {
	db DBTX
}

func NewPostgresAnalyticsRepo(db *sql.DB) *PostgresAnalyticRepo {
	return &PostgresAnalyticRepo{
		db: db,
	}
}