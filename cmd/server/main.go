package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/audit"
	"github.com/JinFuuMugen/ya_go_metrics/internal/compress"
	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/cryptography"
	"github.com/JinFuuMugen/ya_go_metrics/internal/cryptography/rsacrypto"
	"github.com/JinFuuMugen/ya_go_metrics/internal/database"
	"github.com/JinFuuMugen/ya_go_metrics/internal/grpcmetrics"
	"github.com/JinFuuMugen/ya_go_metrics/internal/handlers"
	"github.com/JinFuuMugen/ya_go_metrics/internal/io"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/network"
	pb "github.com/JinFuuMugen/ya_go_metrics/internal/proto"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("cannot create config: %s", err)
	}

	if err := logger.Init(); err != nil {
		log.Fatalf("cannot create logger: %s", err)
	}

	publisher := audit.NewPublisher()

	if cfg.AuditFile != "" {
		fo, err := audit.NewFileObserver(cfg.AuditFile)
		if err != nil {
			log.Fatal(err)
		}
		publisher.Subscribe(fo)
	}

	if cfg.AuditURL != "" {
		publisher.Subscribe(audit.NewHTTPObserver(cfg.AuditURL))
	}

	var db *database.Database
	if cfg.DatabaseDSN != "" {
		db = database.New(cfg.DatabaseDSN)
		if err := db.Connect(); err != nil {
			log.Fatalf("cannot create database connection: %s", err)
		}

		if err := db.Migrate(context.Background()); err != nil {
			log.Fatalf("cannot migrate database: %s", err)
		}
	}

	if err := io.Run(cfg, db); err != nil {
		log.Fatalf("cannot load preload metrics: %s", err)
	}

	st := storage.NewStorage()

	rout := chi.NewRouter()

	if cfg.CryptoKey != "" {
		privateKey, err := rsacrypto.LoadPrivateKey(cfg.CryptoKey)
		if err != nil {
			log.Fatalf("cannot load private key: %s", err)
		}

		rout.Use(rsacrypto.CryptoMiddleware(privateKey))
	}

	rout.Use(compress.GzipMiddleware)

	rout.Mount("/debug", http.DefaultServeMux)

	rout.Get("/", handlers.MainHandler(st))

	rout.Get("/ping", handlers.PingDBHandler(db))

	rout.Route("/updates", func(r chi.Router) {
		r.Use(network.CheckValidSubnetMiddleware(cfg.TrustedSubnet))
		r.Use(cryptography.ValidateHashMiddleware(cfg))
		r.Use(io.GetDumperMiddleware(cfg, db))
		r.Post("/", handlers.UpdateBatchMetricsHandler(st, publisher))
	})

	rout.Route("/update", func(r chi.Router) {
		r.Use(network.CheckValidSubnetMiddleware(cfg.TrustedSubnet))
		r.Use(io.GetDumperMiddleware(cfg, db))
		r.Use(cryptography.ValidateHashMiddleware(cfg))
		r.Post("/", handlers.UpdateMetricsHandler(st, publisher))
		r.Post("/{metric_type}/{metric_name}/{metric_value}", handlers.UpdateMetricsPlainHandler(st, publisher))
	})

	rout.Post("/value/", handlers.GetMetricHandler(st))
	rout.Get("/value/{metric_type}/{metric_name}", handlers.GetMetricPlainHandler(st))

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT,
	)
	defer stop()

	grpcLis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("cannot listen gRPC addr %s: %v", cfg.GRPCAddr, err)
	}

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(network.SubnetUnaryInterceptor(cfg.TrustedSubnet)),
	)

	pb.RegisterMetricsServer(grpcSrv, grpcmetrics.New(st, publisher))

	grpcErrCh := make(chan error, 1)
	go func() {
		logger.Infof("gRPC server listening on %s", cfg.GRPCAddr)
		if err := grpcSrv.Serve(grpcLis); err != nil {
			grpcErrCh <- err
		}
		close(grpcErrCh)
	}()

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: rout,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		logger.Infof("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			logger.Fatalf("cannot start server: %s", err)
		}
	case err := <-grpcErrCh:
		if err != nil {
			logger.Fatalf("cannot start grpc server: %s", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	grpcSrv.GracefulStop()
	_ = grpcLis.Close()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("http server shutdown error: %s", err)
	} else {
		logger.Infof("http server stopped")
	}
}
