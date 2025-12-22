package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

func BenchmarkBatchUpdateHandler(b *testing.B) {
	var metrics []models.Metrics
	for i := 0; i < 100; i++ {
		v := float64(i)
		metrics = append(metrics, models.Metrics{
			ID:    "gauge",
			MType: storage.MetricTypeGauge,
			Value: &v,
		})
	}

	body, _ := json.Marshal(metrics)
	h := UpdateBatchMetricsHandler(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader(body))
		w := httptest.NewRecorder()
		h(w, req)
	}
}
