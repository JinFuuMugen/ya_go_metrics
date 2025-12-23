package fileio

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

// SaveMetricsFile saves all provided metrics to a file in JSON format.
func SaveMetricsFile(filepath string, counters []storage.Counter, gauges []storage.Gauge) error {
	var metrics []models.Metrics

	for _, c := range counters {
		v := c.Value
		metrics = append(metrics, models.Metrics{
			ID:    c.Name,
			MType: storage.MetricTypeCounter,
			Delta: &v,
		})
	}
	for _, g := range gauges {
		v := g.Value
		metrics = append(metrics, models.Metrics{
			ID:    g.Name,
			MType: storage.MetricTypeGauge,
			Value: &v,
		})
	}

	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		return fmt.Errorf("encode metrics: %w", err)
	}
	return nil
}

// LoadMetricsFile loads metrics from a JSON file into in-memory storage.
func LoadMetricsFile(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Warnf("file %s not found. cannot load metrics: %s", filepath, err)
			return nil
		}
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var metrics []models.Metrics
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(&metrics); err != nil {
		return fmt.Errorf("decode metrics: %w", err)
	}

	for _, m := range metrics {
		switch m.MType {
		case storage.MetricTypeCounter:
			if m.Delta == nil {
				return errors.New("counter without delta")
			}
			storage.AddCounter(m.ID, *m.Delta)
		case storage.MetricTypeGauge:
			if m.Value == nil {
				return errors.New("gauge without value")
			}
			storage.SetGauge(m.ID, *m.Value)
		default:
			return fmt.Errorf("unsupported metric type: %s", m.MType)
		}
	}
	return nil
}
