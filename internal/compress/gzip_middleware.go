package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return logger.HandlerLogger(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			if strings.Contains(r.Header.Get("Content-Type"), "application/json") || strings.Contains(r.Header.Get("Content-Type"), "text/html") {
				reader, err := gzip.NewReader(r.Body)
				if err != nil {
					logger.Errorf("cannot create gzip reader: %s", err)
					http.Error(w, fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
					return
				}
				defer reader.Close()
				decodedBody, err := io.ReadAll(reader)
				if err != nil {
					logger.Errorf("cannot decode body: %s", err)
					http.Error(w, fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(decodedBody))
			} else {
				logger.Errorf("invalid content type for gzip encoding")
				http.Error(w, "invalid content type for gzip encoding", http.StatusBadRequest)
				return
			}
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			gzipResponseWriter := &gzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}
			next.ServeHTTP(gzipResponseWriter, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w gzipResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}
