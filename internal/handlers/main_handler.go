package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

// MainHandler handles requests for rendering the main HTML page with all current metric values.
func MainHandler(
	st storage.Storage,
	tmpl *template.Template,
) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {

		if tmpl == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := tmpl.Execute(w, struct {
			Gauges   []storage.Gauge
			Counters []storage.Counter
		}{st.GetGauges(), st.GetCounters()})
		if err != nil {
			logger.Errorf("cannot execute template: %s", err)
			http.Error(w, fmt.Sprintf("cannot execute template: %s", err), http.StatusInternalServerError)
			return
		}
	}
}
