package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestMainHandle(t *testing.T) {
	tests := []struct {
		method     string
		url        string
		name       string
		wantedCode int
	}{
		{
			name:       "wrong method",
			wantedCode: 405,
			method:     http.MethodPost,
			url:        "/",
		},
		{
			name:       "wrong url",
			wantedCode: 404,
			method:     http.MethodGet,
			url:        "/123/",
		},
	}
	logger.Init()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/", MainHandler)
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
