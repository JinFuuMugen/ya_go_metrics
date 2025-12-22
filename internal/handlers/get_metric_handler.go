package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

// GetMetricHandler handles requests for retrieving a single metric value.
// Expects a JSON body containing metric ID and type.
// If the metric is not found, HTTP 404 is returned.
func GetMetricHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var metric models.Metrics
	err := decoder.Decode(&metric)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		logger.Errorf("cannot decode body: %s", err)
		return
	}

	var m storage.Metric

	switch metric.MType {
	case storage.MetricTypeGauge:
		m, err = storage.GetGauge(metric.ID)
		metric.SetValue(m.GetValue().(float64))
	case storage.MetricTypeCounter:
		m, err = storage.GetCounter(metric.ID)
		metric.SetDelta(m.GetValue().(int64))
	default:
		http.Error(w, "unsupported metric type", http.StatusNotImplemented)
		logger.Errorf("unsupported metric type")
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("metric is not found: %s", err), http.StatusNotFound)
		logger.Errorf("metric is not found: %s", err)
		return
	}

	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %s", err), http.StatusInternalServerError)
		logger.Errorf("cannot serialize metric to json: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, "cannot write response", http.StatusInternalServerError)
	}
}
