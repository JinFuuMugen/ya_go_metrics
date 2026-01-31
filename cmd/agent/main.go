package main

import (
	"fmt"
	"log"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/monitors"
	"github.com/JinFuuMugen/ya_go_metrics/internal/sender"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("cannot create config: %s", err)
	}

	err = logger.Init()
	if err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	pollTicker := cfg.PollTicker()
	reportTicker := cfg.ReportTicker()

	str := storage.NewStorage()
	snd := sender.NewSender(*cfg)

	m := monitors.NewRuntimeMonitor(str, snd)
	g := monitors.NewGopsutilMonitor(str, snd)

	rateLimit := cfg.RateLimit
	semaphore := make(chan struct{}, rateLimit)

	rateLimitTicker := time.NewTicker(time.Second / time.Duration(rateLimit))
	defer rateLimitTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			m.CollectRuntimeMetrics()
			if err := g.CollectGopsutil(); err != nil {
				logger.Fatalf("error collecting gopsutil metrics: %s", err)
			}

		case <-reportTicker.C:
			select {
			case semaphore <- struct{}{}:
				go func() {
					err := m.Dump()
					if err != nil {
						logger.Warnf("error dumping metrics: %s", err)
					}
					err = g.Dump()
					if err != nil {
						logger.Warnf("error dumping metrics: %s", err)
					}
					<-semaphore
				}()
			default:
				logger.Warnf("maximum concurrent Dump executions reached, skipping current dump")
			}
		}

		<-rateLimitTicker.C
	}
}
