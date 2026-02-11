package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/cryptography/rsacrypto"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/monitors"
	"github.com/JinFuuMugen/ya_go_metrics/internal/sender"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("cannot create config: %s", err)
	}

	err = logger.Init()
	if err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT,
	)
	defer stop()

	pollTicker := cfg.PollTicker()
	reportTicker := cfg.ReportTicker()
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	str := storage.NewStorage()

	var publicKey *rsa.PublicKey
	if cfg.CryptoKey != "" {
		publicKey, err = rsacrypto.LoadPublicKey(cfg.CryptoKey)
		if err != nil {
			log.Fatalf("cannot load public key: %s", err)
		}
	}

	var snd sender.Sender

	if cfg.GRPCAddr != "" {
		gs, err := sender.NewGRPCSender(*cfg)
		if err != nil {
			log.Fatalf("cannot init grpc sender: %s", err)
		}
		defer gs.Close()
		snd = gs
	} else {
		snd = sender.NewSender(*cfg, publicKey)
	}

	m := monitors.NewRuntimeMonitor(str, snd)
	g := monitors.NewGopsutilMonitor(str, snd)

	rateLimit := cfg.RateLimit
	semaphore := make(chan struct{}, rateLimit)

	rateTicker := time.NewTicker(time.Second / time.Duration(rateLimit))
	defer rateTicker.Stop()

	var wg sync.WaitGroup
	var shuttingDown bool

	waitRateToken := func() bool {
		select {
		case <-rateTicker.C:
			return true
		case <-ctx.Done():
			return false
		}
	}

	finalFlush := func() {
		if err := m.Dump(); err != nil {
			logger.Warnf("final dump runtime metrics error: %s", err)
		}
		if err := g.Dump(); err != nil {
			logger.Warnf("final dump gopsutil metrics error: %s", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			if shuttingDown {
				return
			}
			shuttingDown = true

			logger.Infof("shutdown signal received, stopping agent...")

			pollTicker.Stop()
			reportTicker.Stop()

			wg.Wait()

			finalFlush()

			logger.Infof("agent stopped gracefully")
			return

		case <-pollTicker.C:
			if shuttingDown {
				continue
			}
			if !waitRateToken() {
				continue
			}

			m.CollectRuntimeMetrics()
			if err := g.CollectGopsutil(); err != nil {
				logger.Fatalf("error collecting gopsutil metrics: %s", err)
			}

		case <-reportTicker.C:
			if shuttingDown {
				continue
			}
			if !waitRateToken() {
				continue
			}

			select {
			case semaphore <- struct{}{}:
				wg.Add(1)
				go func() {
					defer wg.Done()
					defer func() {
						<-semaphore
					}()

					if err := m.Dump(); err != nil {
						logger.Warnf("error dumping runtime metrics: %s", err)
					}

					if err := g.Dump(); err != nil {
						logger.Warnf("error dumping gopsutil metrics: %s", err)
					}
				}()
			default:
				logger.Warnf("maximum concurrent Dump executions reached, skipping current dump")
			}
		}
	}
}
