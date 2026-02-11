package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// ServerConfig stores server configuration parameters.
type ServerConfig struct {
	// Addr is the server address in the form host:port.
	Addr string `env:"ADDRESS" json:"address"`

	// StoreInterval defines the interval for saving metrics to persistent storage.
	// A zero value means synchronous saving.
	StoreInterval time.Duration `env:"STORE_INTERVAL" json:"store_interval"`

	// FileStoragePath is the path to the file used for storing metrics.
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`

	// Restore enables or disables restoring metrics on startup.
	Restore bool `env:"RESTORE" json:"restore"`

	// DatabaseDSN is the data source name for connecting to the database.
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn"`

	// Key is an optional key used for SHA256 request signing and verification.
	Key string `env:"KEY" json:"-"`

	// AuditFile defines the file path for audit event logging.
	AuditFile string `env:"AUDIT_FILE" json:"-"`

	// AuditURL defines the HTTP endpoint for sending audit events.
	AuditURL string `env:"AUDIT_URL" json:"-"`

	// CryptoKey is the path to private key file
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`

	// ConfigPath is the path to json config file
	ConfigPath string `env:"CONFIG" json:"-"`
}

// LoadServerConfig loads and initializes the server configuration.
func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{
		Addr:            "localhost:8080",
		StoreInterval:   300 * time.Second,
		FileStoragePath: "tmp/metrics-db.json",
		Restore:         true,
		Key:             "",
		AuditFile:       "",
		AuditURL:        "",
		CryptoKey:       "",
		ConfigPath:      "",
	}

	cfg.ConfigPath = os.Getenv("CONFIG")
	if p, ok := findConfigPath(os.Args[1:]); ok {
		cfg.ConfigPath = p
	}

	if cfg.ConfigPath != "" {
		if err := applyServerJSON(cfg, cfg.ConfigPath); err != nil {
			return nil, fmt.Errorf("cannot apply json config: %w", err)
		}
	}

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "server address")
	flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "metrics store interval(0 to sync)")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "path of storage file")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "boolean to load/not saved values")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database DSN")
	flag.StringVar(&cfg.Key, "k", cfg.Key, "SHA256 key")
	flag.StringVar(&cfg.AuditFile, "audit-file", cfg.AuditFile, "audit file path")
	flag.StringVar(&cfg.AuditURL, "audit-url", cfg.AuditURL, "audit url")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "crypto key filepath")

	flag.Parse()

	if envAddr, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Addr = envAddr
	}

	if envStoreInterval, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		_, err := strconv.Atoi(envStoreInterval)
		if err == nil {
			envStoreInterval = envStoreInterval + "s"
		}
		storeInterval, err := time.ParseDuration(envStoreInterval)
		if err != nil {
			return nil, fmt.Errorf("cannot convert env STORE_INTERVAL to duration value: %w", err)
		}
		cfg.StoreInterval = storeInterval
	}

	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = envFileStoragePath
	}

	if envRestore, ok := os.LookupEnv("RESTORE"); ok {
		restore, err := strconv.ParseBool(envRestore)
		if err != nil {
			return nil, fmt.Errorf("cannot convert env RESTORE to boolean value: %w", err)
		}
		cfg.Restore = restore
	}

	if envDatabaseDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseDSN = envDatabaseDSN
	}

	if envKey, ok := os.LookupEnv("KEY"); ok {
		cfg.Key = envKey
	}

	if envAuditFile, ok := os.LookupEnv("AUDIT_FILE"); ok {
		cfg.AuditFile = envAuditFile
	}

	if envAuditURL, ok := os.LookupEnv("AUDIT_URL"); ok {
		cfg.AuditURL = envAuditURL
	}

	if envCryptoKey, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		cfg.CryptoKey = envCryptoKey
	}

	return cfg, nil
}
