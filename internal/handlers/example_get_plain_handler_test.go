package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/JinFuuMugen/ya_go_metrics/internal/handlers"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

func ExampleGetMetricPlainHandler() {
	st := storage.NewStorage()
	st.AddCounter("requests", 10)

	r := chi.NewRouter()
	r.Get("/value/{metric_type}/{metric_name}", handlers.GetMetricPlainHandler(st))

	req := httptest.NewRequest(http.MethodGet, "/value/counter/requests", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		fmt.Println(rec.Body.String())
	}

	// Output:
	// 10
}
