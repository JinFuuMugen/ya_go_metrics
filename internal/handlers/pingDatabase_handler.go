package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingDBHandler(cfg *config.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dsn := cfg.DatabaseDSN
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			logger.Errorf("can't connect to database: %s", err)
			http.Error(w, "can't connect to database", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = db.PingContext(ctx); err != nil {
			logger.Errorf("error pinging database: %s", err)
			http.Error(w, "error pinging database", http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
