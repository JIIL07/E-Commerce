package database

import (
	"context"
	"database/sql"
)

type Database interface {
	Connect() error
	Close() error
	Ping() error
	GetDB() *sql.DB
	Migrate() error
	Seed() error
}

type Transaction interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
}

type QueryExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type DatabaseConfig interface {
	GetDSN() string
	GetDriver() string
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxLifetime() int
}
