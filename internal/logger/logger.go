package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var log zap.SugaredLogger

func Init() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("cannot initialize zap logger: %w", err)
	}
	sug := *logger.Sugar()
	log = sug
	return nil
}

func Warnf(template string, args ...any) {
	log.Warnf(template, args...)
}

func Fatalf(template string, args ...any) {
	log.Fatalf(template, args...)
}

func Errorf(template string, args ...any) {
	log.Errorf(template, args...)
}

func Infof(template string, args ...any) {
	log.Infof(template, args...)
}

func HandlerLogger(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		responseData := &responseData{
			status: http.StatusOK,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		log.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)

	}
	return logFn
}

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return 0, fmt.Errorf("cannot implement ResponseWriter: %w", err)
	}
	r.responseData.size += size
	return size, nil
}
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
