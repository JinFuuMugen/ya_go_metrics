package handlers

import (
	"encoding/json"
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
	st storage.Storage,
	auditPublisher *audit.Publisher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := readRequestBody(r)
		if err != nil {
			logger.Errorf(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var metrics []models.Metrics
		if err := json.Unmarshal(body, &metrics); err != nil {
			logger.Errorf("cannot unmarshal metrics: %s", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		metricNames := make([]string, 0, len(metrics))

		for _, metric := range metrics {
			switch metric.MType {

			case storage.MetricTypeCounter:
				delta, err := metric.GetDelta()
				if err != nil {
					logger.Errorf("invalid counter metric: %s", err)
					http.Error(w, "bad counter metric", http.StatusBadRequest)
					return
				}
				st.AddCounter(metric.ID, delta)

			case storage.MetricTypeGauge:
				value, err := metric.GetValue()
				if err != nil {
					logger.Errorf("invalid gauge metric: %s", err)
					http.Error(w, "bad gauge metric", http.StatusBadRequest)
					return
				}
				st.SetGauge(metric.ID, value)

			default:
				logger.Errorf("unsupported metric type: %s", metric.MType)
				http.Error(w, "unsupported metric type", http.StatusNotImplemented)
				return
			}

			metricNames = append(metricNames, metric.ID)
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
