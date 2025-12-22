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
	storage.Reset()
	storage.SetGauge("temperature", 36.6)

	reqMetric := models.Metrics{
		ID:    "temperature",
		MType: "gauge",
	}

	body, _ := json.Marshal(reqMetric)

	req := httptest.NewRequest(http.MethodPost, "/value", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handlers.GetMetricHandler(rec, req)

	if rec.Code == http.StatusOK {
		fmt.Println(rec.Body.String())
	}

	// Output:
	// {"id":"temperature","type":"gauge","value":36.6}
}
