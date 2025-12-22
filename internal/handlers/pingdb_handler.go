package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/database"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// PingDBHandler returns an HTTP handler that checks database availability.
func PingDBHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "database is not configured", http.StatusServiceUnavailable)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			logger.Errorf("error pinging database: %s", err)
			http.Error(w, "error pinging database", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
