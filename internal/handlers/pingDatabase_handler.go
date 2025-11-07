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

func PingDBHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := database.DB.PingContext(ctx); err != nil {
			logger.Errorf("error pinging database: %s", err)
			http.Error(w, "error pinging database", http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
