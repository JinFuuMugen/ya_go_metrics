package io

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/JinFuuMugen/ya_go_metrics/internal/database"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func saveMetricsDB(db *database.Database, counters []storage.Counter, gauges []storage.Gauge) error {
	ctx := context.Background()
	query := `
        INSERT INTO metrics (type, value, delta, name)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (type, name) DO UPDATE SET value = $2, delta = $3;
    `

	for _, c := range counters {
		_, err := db.Exec(ctx, query, c.Type, sql.NullFloat64{Valid: false, Float64: 0}, c.Value, c.Name)
		if err != nil {
			return fmt.Errorf("cannot save counter %q: %w", c.Name, err)
		}
	}

	for _, g := range gauges {
		_, err := db.Exec(ctx, query, g.Type, g.Value, sql.NullInt64{Valid: false, Int64: 0}, g.Name)
		if err != nil {
			return fmt.Errorf("cannot save gauge %q: %w", g.Name, err)
		}
	}

	return nil
}

func loadMetricsDB(db *database.Database) error {
	ctx := context.Background()
	rows, err := db.Query(ctx, "SELECT id, type, value, delta FROM metrics")
	if err != nil {
		return fmt.Errorf("cannot read metrics from db: %w", err)
	}
	defer rows.Close()

	var metrics []models.Metrics

	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			return fmt.Errorf("cannot scan values: %w", err)
		}
		metrics = append(metrics, m)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error while iterating metrics: %w", err)
	}

	for _, m := range metrics {
		switch m.MType {
		case storage.MetricTypeCounter:
			if m.Delta != nil {
				storage.AddCounter(m.ID, *m.Delta)
			}
		case storage.MetricTypeGauge:
			if m.Value != nil {
				storage.SetGauge(m.ID, *m.Value)
			}
		default:
			return fmt.Errorf("unsupported metric type: %w", errors.New(m.MType))
		}
	}
	return nil
}
