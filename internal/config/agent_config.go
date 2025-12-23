package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

// Config stores agent configuration parameters.
type Config struct {
	// Addr is the server address in the form host:port.
	Addr string `env:"ADDRESS"`

	// PollInterval defines the interval (in seconds) between metric collection.
	PollInterval int `env:"POLL_INTERVAL"`

	// ReportInterval defines the interval (in seconds) between metric reports.
	ReportInterval int `env:"REPORT_INTERVAL"`

	// Key is an optional key used for SHA256 request signing.
	Key string `env:"KEY"`

	// RateLimit defines the maximum number of outgoing requests.
	RateLimit int `env:"RATE_LIMIT"`
}

// New creates and initializes a Config instace.
func New() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.Addr, `a`, cfg.Addr, `server address`)
	flag.IntVar(&cfg.PollInterval, `p`, cfg.PollInterval, `poll interval`)
	flag.IntVar(&cfg.ReportInterval, `r`, cfg.ReportInterval, `poll interval`)
	flag.StringVar(&cfg.Key, `k`, cfg.Key, `SHA256 key`)
	flag.IntVar(&cfg.RateLimit, `l`, cfg.RateLimit, `requests limit`)
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

	if cfg.RateLimit == 0 {
		cfg.RateLimit = 1
	}

	return cfg, nil
}

// PollTicker returns a ticker that triggers metric collection.
func (cfg *Config) PollTicker() *time.Ticker {
	return time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
}

// ReportTicker return a ticker that triggers metric reporting.
func (cfg *Config) ReportTicker() *time.Ticker {
	return time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
}
