// internal/io/dbio/dbio.go
package dbio

import (
	"context"
	"errors"
	"fmt"

	"github.com/JinFuuMugen/ya_go_metrics/internal/database"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

func SaveMetricsDB(db *database.Database, counters []storage.Counter, gauges []storage.Gauge) error {
	if db == nil {
		return errors.New("db is nil")
	}

	ctx := context.Background()

	for _, c := range counters {
		const q = `
			INSERT INTO metrics (type, name, value, delta)
			VALUES ($1, $2, NULL, $3)
			ON CONFLICT (type, name)
			DO UPDATE SET value = EXCLUDED.value, delta = EXCLUDED.delta;`
		if _, err := db.Exec(ctx, q, storage.MetricTypeCounter, c.Name, c.Value); err != nil {
			return fmt.Errorf("upsert counter %q: %w", c.Name, err)
		}
	}

	for _, g := range gauges {
		const q = `
			INSERT INTO metrics (type, name, value, delta)
			VALUES ($1, $2, $3, NULL)
			ON CONFLICT (type, name)
			DO UPDATE SET value = EXCLUDED.value, delta = EXCLUDED.delta;`
		if _, err := db.Exec(ctx, q, storage.MetricTypeGauge, g.Name, g.Value); err != nil {
			return fmt.Errorf("upsert gauge %q: %w", g.Name, err)
		}
	}

	return nil
}

func LoadMetricsDB(db *database.Database) error {
	if db == nil {
		return errors.New("db is nil")
	}

	ctx := context.Background()

	const q = `SELECT type, name, value, delta FROM metrics`
	rows, err := db.Query(ctx, q)
	if err != nil {
		return fmt.Errorf("select metrics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			typ  string
			name string
			val  *float64
			dlt  *int64
		)
		if err := rows.Scan(&typ, &name, &val, &dlt); err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		switch typ {
		case storage.MetricTypeCounter:
			if dlt == nil {
				return fmt.Errorf("db counter %q without delta", name)
			}
			storage.AddCounter(name, *dlt)
		case storage.MetricTypeGauge:
			if val == nil {
				return fmt.Errorf("db gauge %q without value", name)
			}
			storage.SetGauge(name, *val)
		default:
			return fmt.Errorf("unsupported metric type: %s", typ)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows: %w", err)
	}
	return nil
}
