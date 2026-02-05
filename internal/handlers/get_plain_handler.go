package handlers

import (
	"net/http"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

// GetMetricPlainHandler handles requests for getting a metric in plain text format.
// The handler expects metric type and metric name as URL parameters.
func GetMetricPlainHandler(
	st storage.Storage,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metric_type")
		metricName := chi.URLParam(r, "metric_name")

		w.Header().Set("Content-Type", "text/plain")

		switch metricType {

		case storage.MetricTypeGauge:
			g, err := st.GetGauge(metricName)
			if err != nil {
				http.Error(w, "metric not found", http.StatusNotFound)
				return
			}
			if _, err := w.Write([]byte(g.GetValueString())); err != nil {
				logger.Errorf("cannot write response: %s", err)
			}

		case storage.MetricTypeCounter:
			c, err := st.GetCounter(metricName)
			if err != nil {
				http.Error(w, "metric not found", http.StatusNotFound)
				return
			}
			if _, err := w.Write([]byte(c.GetValueString())); err != nil {
				logger.Errorf("cannot write response: %s", err)
			}

		default:
			logger.Errorf("unsupported metric type: %s", metricType)
			http.Error(w, "unsupported metric type", http.StatusBadRequest)
			return
		}
	}
}
