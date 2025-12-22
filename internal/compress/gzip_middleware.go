package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
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

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
