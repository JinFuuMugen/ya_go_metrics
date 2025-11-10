package io

import (
	"github.com/JinFuuMugen/ya_go_metrics/internal/database"
	"github.com/JinFuuMugen/ya_go_metrics/internal/io/dbio"
	"github.com/JinFuuMugen/ya_go_metrics/internal/io/fileio"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

func saveMetricsFile(path string, counters []storage.Counter, gauges []storage.Gauge) error {
	return fileio.SaveMetricsFile(path, counters, gauges)
}

func loadMetricsFile(path string) error {
	return fileio.LoadMetricsFile(path)
}

func saveMetricsDB(db *database.Database, counters []storage.Counter, gauges []storage.Gauge) error {
	if db == nil {
		return nil
	}
	return dbio.SaveMetricsDB(db, counters, gauges)
}

func loadMetricsDB(db *database.Database) error {
	if db == nil {
		return nil
	}
	return dbio.LoadMetricsDB(db)
}
