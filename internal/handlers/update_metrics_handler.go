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

// UpdateMetricsHandler returns an HTTP handler for metric updates.
// The handler accepts a JSON metric.
func UpdateMetricsHandler(
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

		var metric models.Metrics
		if err := json.Unmarshal(body, &metric); err != nil {
			logger.Errorf("cannot decode body: %s", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		switch metric.MType {

		case storage.MetricTypeCounter:
			delta, err := metric.GetDelta()
			if err != nil {
				logger.Errorf("invalid counter metric: %s", err)
				http.Error(w, "bad counter metric", http.StatusBadRequest)
				return
			}

			st.AddCounter(metric.ID, delta)

			c, err := st.GetCounter(metric.ID)
			if err != nil {
				http.Error(w, "metric not found", http.StatusNotFound)
				return
			}
			metric.SetDelta(c.Value)

		case storage.MetricTypeGauge:
			value, err := metric.GetValue()
			if err != nil {
				logger.Errorf("invalid gauge metric: %s", err)
				http.Error(w, "bad gauge metric", http.StatusBadRequest)
				return
			}

			st.SetGauge(metric.ID, value)
			metric.SetValue(value)

		default:
			logger.Errorf("unsupported metric type: %s", metric.MType)
			http.Error(w, "unsupported metric type", http.StatusNotImplemented)
			return
		}

		if auditPublisher != nil {
			auditPublisher.Publish(models.AuditEvent{
				TS:        time.Now().Unix(),
				Metrics:   []string{metric.ID},
				IPAddress: extractIP(r),
			})
		}

		resp, err := json.Marshal(metric)
		if err != nil {
			logger.Errorf("cannot serialize metric: %s", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(resp); err != nil {
			logger.Errorf("cannot write response: %s", err)
		}
	}
}
