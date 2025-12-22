package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/audit"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

func UpdateMetricsPlainHandler(
	auditPublisher *audit.Publisher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metric_type")
		metricName := chi.URLParam(r, "metric_name")
		metricValue := chi.URLParam(r, "metric_value")

		switch metricType {
		case storage.MetricTypeCounter:
			v, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				logger.Errorf("not a valid metric value: %s", err)
				http.Error(w, fmt.Sprintf("not a valid metric value: %s", err), http.StatusBadRequest)
				return
			}
			storage.AddCounter(metricName, v)
		case storage.MetricTypeGauge:
			v, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				logger.Errorf("not a valid metric value: %s", err)
				http.Error(w, fmt.Sprintf("not a valid metric value: %s", err), http.StatusBadRequest)
				return
			}
			storage.SetGauge(metricName, v)
		default:
			logger.Errorf("unsupported metric type")
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

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
