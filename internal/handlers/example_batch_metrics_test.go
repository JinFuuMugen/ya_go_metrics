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
	storage.Reset()

	handler := handlers.UpdateBatchMetricsHandler(nil)

	body := []byte(`[
		{"id":"requests","type":"counter","delta":5},
		{"id":"temperature","type":"gauge","value":36.6}
	]`)

	req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		fmt.Println("batch metrics updated")
	}

	// Output:
	// batch metrics updated
}
