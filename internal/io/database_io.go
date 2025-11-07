package io

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/JinFuuMugen/ya_go_metrics/internal/database"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func saveMetricsDB(counters []storage.Counter, gauges []storage.Gauge) error {
	for _, c := range counters {
		_, err := database.DB.Exec("INSERT INTO metrics (type, value, delta, id) VALUES ($1, $2, $3, $4) ON CONFLICT (type, id) DO UPDATE SET value = $2, delta = $3;", c.Type, sql.NullFloat64{Valid: false, Float64: 0}, c.Value, c.Name)
		if err != nil {
			return fmt.Errorf("cannot execute query to save counters: %w", err)
		}
	}

	for _, g := range gauges {
		_, err := database.DB.Exec("INSERT INTO metrics (type, value, delta, id) VALUES ($1, $2, $3, $4) ON CONFLICT (type, id) DO UPDATE SET value = $2, delta = $3;", g.Type, g.Value, sql.NullInt64{Valid: false, Int64: 0}, g.Name)
		if err != nil {
			return fmt.Errorf("cannot execute query to save gauges: %w", err)
		}
	}
	return nil
}

func loadMetricsDB() error {
	var metrics []models.Metrics

	rows, err := database.DB.Query("SELECT id, type, value, delta from metrics")
	if err != nil {
		return fmt.Errorf("cannot read metrics from db: %w", err)
	}
	if rows.Err() != nil {
		return fmt.Errorf("cannot read metrics from db: %w", rows.Err())
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			return fmt.Errorf("cannot scan values from db: %w", err)
		}

		metrics = append(metrics, m)
	}

	for _, m := range metrics {
		switch m.MType {
		case storage.MetricTypeCounter:
			storage.AddCounter(m.ID, *m.Delta)
		case storage.MetricTypeGauge:
			storage.SetGauge(m.ID, *m.Value)
		default:
			return fmt.Errorf("cannot opperate metric: %w", errors.New("unsupported metric type"))
		}
	}
	return nil
}
