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

func UpdateMetricsHandler(
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

		var metric models.Metrics
		err = json.Unmarshal(buf.Bytes(), &metric)
		if err != nil {
			logger.Errorf("cannot process body: %s", err)
			http.Error(w, fmt.Sprintf("cannot process body: %s", err), http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case storage.MetricTypeCounter:
			delta, err := metric.GetDelta()
			if err != nil {
				logger.Errorf("cannot get delta: %s", err)
				http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
				return
			}
			storage.AddCounter(metric.ID, delta)
			tmpCounter, _ := storage.GetCounter(metric.ID)
			deltaNew := tmpCounter.GetValue().(int64)
			metric.SetDelta(deltaNew)
		case storage.MetricTypeGauge:
			value, err := metric.GetValue()
			if err != nil {
				logger.Errorf("cannot get value: %s", err)
				http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
				return
			}
			storage.SetGauge(metric.ID, value)
		default:
			logger.Errorf("unsupported metric type")
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

		jsonBytes, err := json.Marshal(metric)
		if err != nil {
			logger.Errorf("cannot serialize metric: %s", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonBytes)
	}
}
