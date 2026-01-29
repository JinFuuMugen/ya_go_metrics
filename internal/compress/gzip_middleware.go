package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
)

var gzipWriterPool = sync.Pool{
	New: func() any {
		return gzip.NewWriter(io.Discard)
	},
}

var gzipReaderPool = sync.Pool{
	New: func() any {
		return new(gzip.Reader)
	},
}

//generate:reset
//go:generate go run ../../cmd/reset/main.go
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

// Write compresses the response data using gzip before wrting it by response writer.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

// GzipMiddleware provides transparent grzip compression and decompression for HTTP requests and responses.
// Incoming requests with Content-Encoding set to "gzip" are decompressed
// before being passed to the next handler.
// If the client supports gzip compression (Accept-Encoding contains "gzip"),
// the response body is compressed before being sent.
func GzipMiddleware(next http.Handler) http.Handler {
	return logger.HandlerLogger(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Encoding") == "gzip" {
			gr := gzipReaderPool.Get().(*gzip.Reader)

			if err := gr.Reset(r.Body); err != nil {
				gzipReaderPool.Put(gr)
				http.Error(w, "invalid gzip body", http.StatusBadRequest)
				return
			}

			r.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: gr,
				Closer: r.Body,
			}

			defer func() {
				gr.Close()
				gzipReaderPool.Put(gr)
			}()
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gw := gzipWriterPool.Get().(*gzip.Writer)
		gw.Reset(w)

		defer func() {
			gw.Close()
			gzipWriterPool.Put(gw)
		}()

		w.Header().Set("Content-Encoding", "gzip")

		grw := &gzipResponseWriter{
			ResponseWriter: w,
			writer:         gw,
		}

		next.ServeHTTP(grw, r)
	})
}
