package handlers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/JinFuuMugen/ya_go_metrics/internal/handlers"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

func ExampleUpdateBatchMetricsHandler() {
	st := storage.NewStorage()

	handler := handlers.UpdateBatchMetricsHandler(st, nil)

	body := []byte(`[
		{"id":"requests","type":"counter","delta":5},
		{"id":"temperature","type":"gauge","value":36.6}
	]`)

	req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	counter, _ := st.GetCounter("requests")
	gauge, _ := st.GetGauge("temperature")

	fmt.Println(counter.Value)
	fmt.Println(gauge.Value)

	// Output:
	// 5
	// 36.6
}
