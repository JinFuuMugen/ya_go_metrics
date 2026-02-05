package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func findConfigPath(args []string) (string, bool) {
	for i := range args {
		a := args[i]
		if a == "-c" || a == "-config" {
			if i+1 < len(args) {
				return args[i+1], true
			}
			return "", false
		}
		if after, ok := strings.CutPrefix(a, "-c="); ok {
			return after, true
		}
		if after, ok := strings.CutPrefix(a, "-config="); ok {
			return after, true
		}
	}
	return "", false
}

func applyServerJSON(cfg *ServerConfig, path string) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cannot parse config json: %w", err)
	}

	if v, ok := raw["address"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil && s != "" {
			cfg.Addr = s
		}
	}

	if v, ok := raw["store_interval"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			return fmt.Errorf("invalid store_interval: %w", err)
		}
		if s != "" {
			d, err := time.ParseDuration(s)
			if err != nil {
				return fmt.Errorf("invalid store_interval: %w", err)
			}
			cfg.StoreInterval = d
		}
	}

	if v, ok := raw["store_file"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil && s != "" {
			cfg.FileStoragePath = s
		}
	}

	if v, ok := raw["restore"]; ok {
		var b bool
		if err := json.Unmarshal(v, &b); err != nil {
			return fmt.Errorf("invalid restore: %w", err)
		}
		cfg.Restore = b
	}

	if v, ok := raw["database_dsn"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil && s != "" {
			cfg.DatabaseDSN = s
		}
	}

	if v, ok := raw["crypto_key"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil && s != "" {
			cfg.CryptoKey = s
		}
	}

	return nil
}

func applyAgentJSON(cfg *AgentConfig, path string) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cannot parse config json: %w", err)
	}

	if v, ok := raw["address"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil && s != "" {
			cfg.Addr = s
		}
	}

	if v, ok := raw["poll_interval"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			return fmt.Errorf("invalid poll_interval: %w", err)
		}
		if s != "" {
			d, err := time.ParseDuration(s)
			if err != nil {
				return fmt.Errorf("invalid poll_interval: %w", err)
			}
			cfg.PollInterval = int(d.Seconds())
		}
	}

	if v, ok := raw["report_interval"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			return fmt.Errorf("invalid report_interval: %w", err)
		}
		if s != "" {
			d, err := time.ParseDuration(s)
			if err != nil {
				return fmt.Errorf("invalid report_interval: %w", err)
			}
			cfg.ReportInterval = int(d.Seconds())
		}
	}

	if v, ok := raw["crypto_key"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil && s != "" {
			cfg.CryptoKey = s
		}
	}

	return nil
}
