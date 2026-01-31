package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Database provides a wrapper around sql.DB for working with PostgreSQL.
type Database struct {
	db  *sql.DB
	dsn string
}

// New creates a new Database instance using the provided DSN.
func New(dsn string) *Database {
	return &Database{dsn: dsn}
}

// Connect establishes a connection to the database.
func (d *Database) Connect() error {
	if d.db != nil {
		return nil
	}
	db, err := sql.Open("pgx", d.dsn)
	if err != nil {
		return fmt.Errorf("cannot connect to database: %w", err)
	}
	d.db = db
	return nil
}

// Migrate initializes the database schema if it does not exist.
func (d *Database) Migrate(ctx context.Context) error {
	if d.db == nil {
		return fmt.Errorf("database is not connected")
	}
	_, err := d.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS metrics (id TEXT PRIMARY KEY DEFAULT gen_random_uuid(), type TEXT NOT NULL,  name TEXT NOT NULL, value DOUBLE PRECISION, delta INT); CREATE UNIQUE INDEX IF NOT EXISTS idx_metrics_type_name ON metrics USING btree (type, name);")

	if err != nil {
		return fmt.Errorf("cannot create table: %w", err)
	}

	return nil
}

// Exec executes a query without returning any rows.
func (d *Database) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if d.db == nil {
		return nil, fmt.Errorf("database is not connected")
	}
	return d.db.ExecContext(ctx, query, args...)
}

// Ping verifies that the database connection is alive.
func (d *Database) Ping(ctx context.Context) error {
	if d.db == nil {
		return fmt.Errorf("database is not connected")
	}
	return d.db.PingContext(ctx)
}

// Query executes a query that returns rows.
func (d *Database) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if d.db == nil {
		return nil, fmt.Errorf("database is not connected")
	}
	return d.db.QueryContext(ctx, query, args...)
}

// Close closes the database connection.
func (d *Database) Close() error {
	if d.db == nil {
		return nil
	}
	err := d.db.Close()
	d.db = nil
	return err
}
