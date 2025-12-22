package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/audit"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

// UpdateBatchMetricsHandler returns an HTTP handler for batch metric updates.
// The handler accepts a JSON array of metrics.
func UpdateBatchMetricsHandler(
	auditPublisher *audit.Publisher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			logger.Errorf("cannot read request body: %s", err)
			http.Error(w, fmt.Sprintf("cannot read request body: %s", err), http.StatusBadRequest)
			return
		}

		var metrics []models.Metrics
		err = json.Unmarshal(buf.Bytes(), &metrics)
		if err != nil {
			logger.Errorf("cannot process body: %s", err)
			http.Error(w, fmt.Sprintf("cannot process body: %s", err), http.StatusBadRequest)
			return
		}

		metricNames := make([]string, 0, len(metrics))

		for _, metric := range metrics {
			switch metric.MType {

			case storage.MetricTypeCounter:
				delta, err := metric.GetDelta()
				if err != nil {
					logger.Errorf("cannot get counter delta: %s", err)
					http.Error(w, "bad counter metric", http.StatusBadRequest)
					return
				}
				storage.AddCounter(metric.ID, delta)
				metricNames = append(metricNames, metric.ID)

			case storage.MetricTypeGauge:
				value, err := metric.GetValue()
				if err != nil {
					logger.Errorf("cannot get gauge value: %s", err)
					http.Error(w, "bad gauge metric", http.StatusBadRequest)
					return
				}
				storage.SetGauge(metric.ID, value)
				metricNames = append(metricNames, metric.ID)

			default:
				logger.Errorf("unsupported metric type: %s", metric.MType)
				http.Error(w, "unsupported metric type", http.StatusNotImplemented)
				return
			}
		}

		if auditPublisher != nil {
			auditPublisher.Publish(models.AuditEvent{
				TS:        time.Now().Unix(),
				Metrics:   metricNames,
				IPAddress: extractIP(r),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
