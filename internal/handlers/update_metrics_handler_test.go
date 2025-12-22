package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricsHandle(t *testing.T) {
	testGauge := `{"id":"Some", "type":"gauge", "value":124.24}`
	testCounter := `{"id":"Some", "type":"counter", "delta":124}`
	testWrongMetric := `{"id":"Some", "type":"qwert", "delta":124}`
	testWrongValue := `{"id":"Some", "type":"counter", "delta":124.123}`
	testValue := 124.24
	testDelta := int64(124)
	testDoubleDelta := int64(124 * 2)

	tests := []struct {
		method     string
		name       string
		url        string
		wantedCode int
		body       string
		wantedBody models.Metrics
	}{
		{
			name:       "positive gauge post",
			wantedCode: 200,
			method:     http.MethodPost,
			url:        "/update/",
			body:       testGauge,
			wantedBody: models.Metrics{
				ID:    "Some",
				MType: "gauge",
				Delta: nil,
				Value: &testValue,
			},
		},
		{
			name:       "positive counter post",
			wantedCode: 200,
			method:     http.MethodPost,
			url:        "/update/",
			body:       testCounter,
			wantedBody: models.Metrics{
				ID:    "Some",
				MType: "counter",
				Delta: &testDelta,
				Value: nil,
			},
		},
		{
			name:       "positive update existing counter post",
			wantedCode: 200,
			method:     http.MethodPost,
			url:        "/update/",
			body:       testCounter,
			wantedBody: models.Metrics{
				ID:    "Some",
				MType: "counter",
				Delta: &testDoubleDelta,
				Value: nil,
			},
		},
		{
			name:       "wrong method",
			wantedCode: 405,
			method:     http.MethodGet,
			url:        "/update/",
			body:       "",
		},
		{
			name:       "wrong url",
			wantedCode: 404,
			method:     http.MethodPost,
			url:        "/updat",
		},
		{
			name:       "wrong metric",
			wantedCode: 501,
			method:     http.MethodPost,
			url:        "/update/",
			body:       testWrongMetric,
		},
		{
			name:       "bad metric value",
			wantedCode: 400,
			method:     http.MethodPost,
			url:        "/update/",
			body:       testWrongValue,
		},
	}
	logger.Init()
	storage.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/update/", UpdateMetricsHandler(nil))
			req, err := http.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantedCode, rr.Code)
			if tt.wantedCode == 200 {
				var data models.Metrics
				json.Unmarshal(rr.Body.Bytes(), &data)
				assert.Equal(t, tt.wantedBody, data)
			}
		})
	}
}
