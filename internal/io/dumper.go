package io

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/database"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

// Run initializes metric persistence according to the server configuration.
func Run(cfg *config.ServerConfig, db *database.Database) error {
	if cfg.FileStoragePath != "" {

		if cfg.Restore {
			if cfg.DatabaseDSN == "" {
				err := loadMetricsFile(cfg.FileStoragePath)
				if err != nil {
					return fmt.Errorf("cannot read metrics from file: %w", err)
				}
			} else {
				err := loadMetricsDB(db)
				if err != nil {
					return fmt.Errorf("cannot read metrics from database: %w", err)
				}
			}
		}
	}
	if cfg.StoreInterval > 0 {
		go runDumper(cfg, db)
	}
	return nil
}

func runDumper(cfg *config.ServerConfig, db *database.Database) {
	storeTicker := time.NewTicker(cfg.StoreInterval)
	for range storeTicker.C {
		if cfg.DatabaseDSN != "" {
			if db == nil {
				logger.Errorf("no database handle available to save metrics")
				continue
			}

			err := saveMetricsDB(db, storage.GetCounters(), storage.GetGauges())
			if err != nil {
				logger.Errorf("cannot save metrics into db: %s", err)
			}
		} else {
			err := saveMetricsFile(cfg.FileStoragePath, storage.GetCounters(), storage.GetGauges())
			if err != nil {
				logger.Fatalf("cannot save metrics into file: %s", err)
			}
		}

	}
}

// GetDumperMiddleware returns an HTTP middleware that triggers metric persistence after request handling.
// When StoreInterval is set to zero or less, metrics are saved synchronously after each request.
// Metrics are saved either to file or database depending on the server configuration.
func GetDumperMiddleware(cfg *config.ServerConfig, db *database.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			if cfg.StoreInterval <= 0 {
				if cfg.DatabaseDSN == "" {
					if err := saveMetricsFile(cfg.FileStoragePath, storage.GetCounters(), storage.GetGauges()); err != nil {
						logger.Errorf("cannot write metrics into file: %s", err)
					}
				} else {
					if db == nil {
						logger.Errorf("no database handle available to save metrics")
						return
					}
					if err := saveMetricsDB(db, storage.GetCounters(), storage.GetGauges()); err != nil {
						logger.Errorf("cannot write metrics into db: %s", err)
					}
				}
			}
		})
	}
}
