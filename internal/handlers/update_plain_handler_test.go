package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricsHandlePlain(t *testing.T) {

	tests := []struct {
		method     string
		name       string
		url        string
		wantedCode int
	}{
		{
			name:       "positive gauge post",
			wantedCode: 200,
			method:     http.MethodPost,
			url:        "/update/gauge/someValue/120.414",
		},
		{
			name:       "positive counter post",
			wantedCode: 200,
			method:     http.MethodPost,
			url:        "/update/counter/someValue/120",
		},
		{
			name:       "wrong method",
			wantedCode: 405,
			method:     http.MethodGet,
			url:        "/update/counter/someValue/120",
		},
		{
			name:       "wrong url",
			wantedCode: 404,
			method:     http.MethodPost,
			url:        "/update/",
		},
		{
			name:       "wrong metric",
			wantedCode: 501,
			method:     http.MethodPost,
			url:        "/update/metr/someValue/900.009",
		},
		{
			name:       "bad metric value",
			wantedCode: 400,
			method:     http.MethodPost,
			url:        "/update/counter/someValue/120.321",
		},
	}
	logger.Init()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/update/{metric_type}/{metric_name}/{metric_value}", UpdateMetricsPlainHandler(nil))

			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantedCode, rr.Code)
		})
	}
}
