package io

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestSaveMetrics(t *testing.T) {
	counters := []storage.Counter{
		{Name: "testCNT1", Type: "counter", Value: 10},
		{Name: "testCNT2", Type: "counter", Value: 2},
	}
	gauges := []storage.Gauge{
		{Name: "testGAU1", Type: "gauge", Value: 123.123},
		{Name: "testGAU2", Type: "gauge", Value: 999999.1313},
	}

	filepath := `test_metrics.json`

	err := saveMetricsFile(filepath, counters, gauges)
	if err != nil {
		t.Fatalf("SaveMetrics failed with error: %v", err)
	}

	// read file content and deserialize json
	file, err := os.Open(filepath)
	if err != nil {
		t.Fatalf("Failed to open test metrics file: %v", err)
	}
	defer file.Close()

	var metrics []models.Metrics
	jsonDecoder := json.NewDecoder(file)
	if err := jsonDecoder.Decode(&metrics); err != nil {
		t.Fatalf("Failed to decode JSON from test metrics file: %v", err)
	}

	testCNT1 := int64(10)
	testCNT2 := int64(2)
	testGAU1 := 123.123
	testGAU2 := 999999.1313

	expectedMetrics := []models.Metrics{
		{ID: "testCNT1", MType: "counter", Delta: &testCNT1, Value: nil},
		{ID: "testCNT2", MType: "counter", Delta: &testCNT2, Value: nil},
		{ID: "testGAU1", MType: "gauge", Delta: nil, Value: &testGAU1},
		{ID: "testGAU2", MType: "gauge", Delta: nil, Value: &testGAU2},
	}
	assert.ElementsMatch(t, expectedMetrics, metrics)
}
