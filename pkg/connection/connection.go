package connection

import (
	"GraphQL/configs"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // import the pq driver
	"github.com/pressly/goose/v3"
)

type DBops interface {
	// database queries
	GetPool() *sqlx.DB
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	StartTransaction() (*sqlx.Tx, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Close() error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type Database struct {
	db *sqlx.DB
}

func (s *Database) Close() error {
	//if err := goose.Down(s.db.DB, "./internal/migrations"); err != nil {
	//	fmt.Printf("goose migration down failed: %v", err)
	//}
	return s.db.Close()
}
func (s *Database) StartTransaction() (*sqlx.Tx, error) {
	return s.db.Beginx()
}
func (s *Database) GetPool() *sqlx.DB {
	return s.db
}
func (s *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}
func (s *Database) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.GetContext(ctx, dest, query, args...)
}
func (s *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}
func (s *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}
func (s *Database) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.SelectContext(ctx, dest, query, args...)
}

func GenerateDsn(cfgs *configs.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfgs.Postgres.Host, cfgs.Postgres.Port, cfgs.Postgres.User, cfgs.Postgres.Password, cfgs.Postgres.DBName)
}

// cfgs database configuration from env file.
func NewDB(ctx context.Context, cfgs *configs.Config) (*Database, error) {
	db, err := sqlx.Connect("postgres", GenerateDsn(cfgs))
	if err != nil {
		return nil, fmt.Errorf("could not create connection pool: %w", err)
	}

	if err := goose.Up(db.DB, "./internal/migrations"); err != nil {
		return nil, fmt.Errorf("goose migration up failed: %v", err)
	}

	return &Database{db: db}, nil
}
