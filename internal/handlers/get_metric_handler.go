package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

// GetMetricHandler handles requests for retrieving a single metric value.
// Expects a JSON body containing metric ID and type.
// If the metric is not found, HTTP 404 is returned.
func GetMetricHandler(
	st storage.Storage,
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

		case storage.MetricTypeGauge:
			g, err := st.GetGauge(metric.ID)
			if err != nil {
				http.Error(w, "metric not found", http.StatusNotFound)
				return
			}
			metric.SetValue(g.Value)

		case storage.MetricTypeCounter:
			c, err := st.GetCounter(metric.ID)
			if err != nil {
				http.Error(w, "metric not found", http.StatusNotFound)
				return
			}
			metric.SetDelta(c.Value)

		default:
			logger.Errorf("unsupported metric type: %s", metric.MType)
			http.Error(w, "unsupported metric type", http.StatusNotImplemented)
			return
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
