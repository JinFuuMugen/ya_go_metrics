package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/JinFuuMugen/ya_go_metrics/internal/handlers"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

func ExampleGetMetricHandler() {
	st := storage.NewStorage()
	st.SetGauge("temperature", 36.6)

	handler := handlers.GetMetricHandler(st)

	reqMetric := models.Metrics{
		ID:    "temperature",
		MType: storage.MetricTypeGauge,
	}

	body, _ := json.Marshal(reqMetric)

	req := httptest.NewRequest(http.MethodPost, "/value", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	fmt.Println(rec.Body.String())

	// Output:
	// {"id":"temperature","type":"gauge","value":36.6}
}
