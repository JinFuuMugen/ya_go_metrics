package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
)

// AgentConfig stores agent configuration parameters.
type AgentConfig struct {
	// Addr is the server address in the form host:port.
	Addr string `env:"ADDRESS" json:"address"`

	// PollInterval defines the interval (in seconds) between metric collection.
	PollInterval int `env:"POLL_INTERVAL" json:"poll_interval"`

	// ReportInterval defines the interval (in seconds) between metric reports.
	ReportInterval int `env:"REPORT_INTERVAL" json:"report_interval"`

	// Key is an optional key used for SHA256 request signing.
	Key string `env:"KEY" json:"-"`

	// RateLimit defines the maximum number of outgoing requests.
	RateLimit int `env:"RATE_LIMIT" json:"-"`

	// CryptoKey is the path to private key file
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`

	// ConfigPath is the path to json config file
	ConfigPath string `env:"CONFIG" json:"-"`
}

// New creates and initializes a AgentConfig instace.
func LoadAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}

	cfg.ConfigPath = os.Getenv("CONFIG")
	if p, ok := findConfigPath(os.Args[1:]); ok {
		cfg.ConfigPath = p
	}

	if cfg.ConfigPath != "" {
		if err := applyAgentJSON(cfg, cfg.ConfigPath); err != nil {
			return nil, fmt.Errorf("cannot apply json config: %w", err)
		}
	}

	flag.StringVar(&cfg.Addr, `a`, cfg.Addr, `server address`)
	flag.IntVar(&cfg.PollInterval, `p`, cfg.PollInterval, `poll interval`)
	flag.IntVar(&cfg.ReportInterval, `r`, cfg.ReportInterval, `poll interval`)
	flag.StringVar(&cfg.Key, `k`, cfg.Key, `SHA256 key`)
	flag.IntVar(&cfg.RateLimit, `l`, cfg.RateLimit, `requests limit`)
	flag.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "crypto key filepath")
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

	if envCryptoKey, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		cfg.CryptoKey = envCryptoKey
	}

	return cfg, nil
}

// PollTicker returns a ticker that triggers metric collection.
func (cfg *AgentConfig) PollTicker() *time.Ticker {
	return time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
}

// ReportTicker return a ticker that triggers metric reporting.
func (cfg *AgentConfig) ReportTicker() *time.Ticker {
	return time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
}
