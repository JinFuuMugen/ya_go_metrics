package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Addr            string        `env:"ADDRESS"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
	DatabaseDSN     string        `env:"DATABASE_DSN"`
	Key             string        `env:"KEY"`
	AuditFile       string        `env:"AUDIT_FILE"`
	AuditURL        string        `env:"AUDIT_URL"`
}

func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{
		Addr:            "localhost:8080",
		StoreInterval:   300 * time.Second,
		FileStoragePath: "tmp/metrics-db.json",
		Restore:         true,
		Key:             "",
		AuditFile:       "",
		AuditURL:        "",
	}

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "server address")
	flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "metrics store interval(0 to sync)")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "path of storage file")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "boolean to load/not saved values")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database DSN")
	flag.StringVar(&cfg.Key, "k", cfg.Key, "SHA256 key")
	flag.StringVar(&cfg.AuditFile, "audit-file", cfg.AuditFile, "audit file path")
	flag.StringVar(&cfg.AuditURL, "audit-url", cfg.AuditURL, "audit url")
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

	return cfg, nil
}
