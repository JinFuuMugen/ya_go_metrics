package fileio

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

func saveMetrics(filepath string, counters []storage.Counter, gauges []storage.Gauge) error {
	var metrics []models.Metrics

	for _, c := range counters {
		cDelta := c.GetValue().(int64)
		metrics = append(metrics, models.Metrics{
			ID:    c.GetName(),
			MType: c.GetType(),
			Delta: &cDelta,
			Value: nil,
		})
	}
	for _, g := range gauges {
		gValue := g.GetValue().(float64)
		metrics = append(metrics, models.Metrics{
			ID:    g.GetName(),
			MType: g.GetType(),
			Delta: nil,
			Value: &gValue,
		})
	}

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("cannot serialize metric to json: %w", err)
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file to save metrics: %w", err)
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("cannot truncate file: %w", err)
	}

	if _, err = file.Write(jsonData); err != nil {
		return fmt.Errorf("cannot write json to file: %w", err)
	}
	return nil
}

func loadMetrics(filepath string) error {
	var metrics []models.Metrics

	file, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file to load metrics: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil
	}

	fileData := scanner.Bytes()
	err = json.Unmarshal(fileData, &metrics)
	if err != nil {
		return fmt.Errorf("cannot deserialize file data: %w", err)
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

func Run(cfg *config.ServerConfig) error {
	if cfg.FileStoragePath != "" {

		if cfg.Restore {
			err := loadMetrics(cfg.FileStoragePath)
			if err != nil {
				return fmt.Errorf("cannot read metrics: %w", err)
			}
		}
	}
	if cfg.StoreInterval > 0 {
		go runDumper(cfg)
	}
	return nil
}

func runDumper(cfg *config.ServerConfig) {
	storeTicker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
	for range storeTicker.C {
		err := saveMetrics(cfg.FileStoragePath, storage.GetCounters(), storage.GetGauges())
		if err != nil {
			logger.Fatalf("cannot save metrics: %s", err)
		}
	}
}

func GetDumperMiddleware(cfg *config.ServerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)

			if cfg.StoreInterval <= 0 {
				err := saveMetrics(cfg.FileStoragePath, storage.GetCounters(), storage.GetGauges())
				if err != nil {
					logger.Fatalf("cannot write metrics: %w", err)
				}
			}
		})
	}
}
