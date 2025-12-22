package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

// MainHandler handles requests for rendering the main HTML page with all current metric values.
func MainHandler(w http.ResponseWriter, _ *http.Request) {

	tmpl, err := template.ParseFiles("internal/static/index.html")
	if err != nil {
		logger.Errorf("cannot parse template: %s", err)
		http.Error(w, fmt.Sprintf("cannot parse template: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, struct {
		Gauges   []storage.Gauge
		Counters []storage.Counter
	}{storage.GetGauges(), storage.GetCounters()})
	if err != nil {
		logger.Errorf("cannot execute template: %s", err)
		http.Error(w, fmt.Sprintf("cannot execute template: %s", err), http.StatusInternalServerError)
		return
	}
}
