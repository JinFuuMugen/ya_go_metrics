package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDatabase(config string) error {
	var err error
	DB, err = sql.Open("pgx", config)
	if err != nil {
		return fmt.Errorf("cannot connect to database: %w", err)
	}

	ctx := context.Background()

	_, err = DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS metrics (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), type TEXT NOT NULL, value DOUBLE PRECISION, delta INT, name TEXT NOT NULL); CREATE UNIQUE INDEX IF NOT EXISTS idx_metrics_type_name ON metrics (type, name);")

	if err != nil {
		return fmt.Errorf("cannot create table: %w", err)
	}

	return nil
}

func CloseDatabase() error {
	err := DB.Close()
	if err != nil {
		return fmt.Errorf("cannot close connection to database: %w", err)
	}
	return nil
}
