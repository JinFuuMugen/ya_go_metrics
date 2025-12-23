package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/audit"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

// UpdateMetricsPlainHandler returns an HTTP handler that updates a single metric using plain-text URL parameters.
func UpdateMetricsPlainHandler(
	st storage.Storage,
	auditPublisher *audit.Publisher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metric_type")
		metricName := chi.URLParam(r, "metric_name")
		metricValue := chi.URLParam(r, "metric_value")

		switch metricType {

		case storage.MetricTypeCounter:
			delta, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				logger.Errorf("invalid counter value: %s", err)
				http.Error(w, "invalid counter value", http.StatusBadRequest)
				return
			}
			st.AddCounter(metricName, delta)

		case storage.MetricTypeGauge:
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				logger.Errorf("invalid gauge value: %s", err)
				http.Error(w, "invalid gauge value", http.StatusBadRequest)
				return
			}
			st.SetGauge(metricName, value)

		default:
			logger.Errorf("unsupported metric type: %s", metricType)
			http.Error(w, "unsupported metric type", http.StatusNotImplemented)
			return
		}

		if auditPublisher != nil {
			auditPublisher.Publish(models.AuditEvent{
				TS:        time.Now().Unix(),
				Metrics:   []string{metricName},
				IPAddress: extractIP(r),
			})
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
