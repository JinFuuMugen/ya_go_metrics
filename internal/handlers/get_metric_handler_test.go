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

func TestGetMetricHandle(t *testing.T) {
	testGauge := `{"id":"GetTestGauge", "type":"gauge"}`
	testCounter := `{"id":"GetTestCounter", "type":"counter"}`
	testWrongMetric := `{"id":"Some", "type":"qwert"}`
	testValue := 123.123
	testDelta := int64(123)

	tests := []struct {
		method     string
		name       string
		url        string
		wantedCode int
		body       string
		wantedBody models.Metrics
	}{
		{
			name:       "positive gauge get",
			wantedCode: 200,
			method:     http.MethodPost,
			url:        "/value/",
			body:       testGauge,
			wantedBody: models.Metrics{
				ID:    "GetTestGauge",
				MType: "gauge",
				Delta: nil,
				Value: &testValue,
			},
		},
		{
			name:       "positive counter get",
			wantedCode: 200,
			method:     http.MethodPost,
			url:        "/value/",
			body:       testCounter,
			wantedBody: models.Metrics{
				ID:    "GetTestCounter",
				MType: "counter",
				Delta: &testDelta,
				Value: nil,
			},
		},
		{
			name:       "wrong method",
			wantedCode: 405,
			method:     http.MethodGet,
			url:        "/value/",
			body:       testGauge,
		},
		{
			name:       "wrong url",
			wantedCode: 404,
			method:     http.MethodPost,
			url:        "/valu/",
			body:       testGauge,
		},
		{
			name:       "wrong metric",
			wantedCode: 501,
			method:     http.MethodPost,
			url:        "/value/",
			body:       testWrongMetric,
		},
	}
	logger.Init()
	st := storage.NewStorage()
	st.SetGauge("GetTestGauge", testValue)
	st.AddCounter("GetTestCounter", testDelta)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/value/", GetMetricHandler(st))
			req, err := http.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
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
