package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/JinFuuMugen/ya_go_metrics/internal/handlers"
)

func ExamplePingDBHandler() {
	handler := handlers.PingDBHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code == http.StatusServiceUnavailable {
		fmt.Println("database is not configured")
	}

	// Output:
	// database is not configured
}
