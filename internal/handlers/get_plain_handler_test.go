package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetMetricPlainHandle(t *testing.T) {
	tests := []struct {
		method      string
		name        string
		url         string
		wantedCode  int
		wantedValue string
	}{
		{
			name:        "positive gauge get",
			wantedCode:  http.StatusOK,
			method:      http.MethodGet,
			url:         "/value/gauge/someG",
			wantedValue: "123.123",
		},
		{
			name:        "positive counter get",
			wantedCode:  http.StatusOK,
			method:      http.MethodGet,
			url:         "/value/counter/someC",
			wantedValue: "123",
		},
		{
			name:       "wrong method",
			wantedCode: http.StatusMethodNotAllowed,
			method:     http.MethodPost,
			url:        "/value/counter/someValue",
		},
		{
			name:       "wrong url",
			wantedCode: http.StatusNotFound,
			method:     http.MethodGet,
			url:        "/updat/gauge/some",
		},
		{
			name:       "wrong metric",
			wantedCode: http.StatusBadRequest,
			method:     http.MethodGet,
			url:        "/value/metr/someValue",
		},
	}
	logger.Init()

	st := storage.NewStorage()
	st.SetGauge("someG", 123.123)
	st.AddCounter("someC", 123)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/value/{metric_type}/{metric_name}", GetMetricPlainHandler(st))
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantedCode, rr.Code)
			if tt.wantedCode == 200 {
				assert.Equal(t, tt.wantedValue, rr.Body.String())
			}
		})
	}
}
