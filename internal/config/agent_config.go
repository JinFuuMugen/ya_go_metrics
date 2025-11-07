package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"time"
)

type Config struct {
	Addr           string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func New() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.Addr, `a`, cfg.Addr, `server address`)
	flag.IntVar(&cfg.PollInterval, `p`, cfg.PollInterval, `poll interval`)
	flag.IntVar(&cfg.ReportInterval, `r`, cfg.ReportInterval, `poll interval`)
	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("cannot read env config: %w", err)
	}

	if cfg.Addr == "" {
		cfg.Addr = "localhost:8080"
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = 2
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = 10
	}
	return cfg, nil
}

func (cfg *Config) PollTicker() *time.Ticker {
	return time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
}

func (cfg *Config) ReportTicker() *time.Ticker {
	return time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
}
