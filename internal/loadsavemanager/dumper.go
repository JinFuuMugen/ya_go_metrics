package loadsavemanager

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

func Run(cfg *config.ServerConfig) error {
	if cfg.FileStoragePath != "" {

		if cfg.Restore {
			if cfg.DatabaseDSN == "" {
				err := loadMetricsFile(cfg.FileStoragePath)
				if err != nil {
					return fmt.Errorf("cannot read metrics from file: %w", err)
				}
			} else {
				err := loadMetricsDB()
				if err != nil {
					return fmt.Errorf("cannot read metrics from database: %w", err)
				}
			}
		}
	}
	if cfg.StoreInterval > 0 {
		go runDumper(cfg)
	}
	return nil
}

func runDumper(cfg *config.ServerConfig) {
	storeTicker := time.NewTicker(cfg.StoreInterval)
	for range storeTicker.C {
		if cfg.DatabaseDSN != "" {
			err := saveMetricsDB(storage.GetCounters(), storage.GetGauges())
			if err != nil {
				logger.Fatalf("cannot save metrics into db: %s", err)
			}
		} else {
			err := saveMetricsFile(cfg.FileStoragePath, storage.GetCounters(), storage.GetGauges())
			if err != nil {
				logger.Fatalf("cannot save metrics into file: %s", err)
			}
		}
	}
}

func GetDumperMiddleware(cfg *config.ServerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)

			if cfg.StoreInterval <= 0 {

				if cfg.DatabaseDSN == "" {
					err := saveMetricsFile(cfg.FileStoragePath, storage.GetCounters(), storage.GetGauges())
					if err != nil {
						logger.Fatalf("cannot write metrics into file: %s", err)
					}
				} else {
					err := saveMetricsDB(storage.GetCounters(), storage.GetGauges())
					if err != nil {
						logger.Fatalf("cannot write metrics into db: %s", err)
					}
				}
			}
		})
	}
}
