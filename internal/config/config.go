package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the global Tigclaw configuration
type Config struct {
	// ListenAddr is the address the gateway listens on (e.g., ":443" or ":9000")
	ListenAddr string `json:"listen_addr"`
	// UpstreamAddr is the OpenClaw backend address (e.g., "http://127.0.0.1:3001")
	UpstreamAddr string `json:"upstream_addr"`
	// DataDir is the directory where Tigclaw stores its database and keys
	DataDir string `json:"data_dir"`
	// StrictMode rejects requests with non-tigclaw keys if true
	StrictMode bool `json:"strict_mode"`
	// RateLimit is the maximum requests per second per IP (0 = unlimited)
	RateLimit int `json:"rate_limit"`
}

// DefaultConfig returns sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		ListenAddr:   ":9000",
		UpstreamAddr: "http://127.0.0.1:3001",
		DataDir:      filepath.Join(homeDir, ".tigclaw"),
		StrictMode:   true,
		RateLimit:    60,
	}
}

// ConfigPath returns the path to the config file
func ConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".tigclaw", "config.json")
}

// Load reads configuration from disk, falling back to defaults
func Load() (*Config, error) {
	cfg := DefaultConfig()
	path := ConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// First run — create default config
			if saveErr := cfg.Save(); saveErr != nil {
				return nil, fmt.Errorf("failed to create default config: %w", saveErr)
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}

// Save writes the current configuration to disk
func (c *Config) Save() error {
	if err := os.MkdirAll(filepath.Dir(ConfigPath()), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0600)
}

// DBPath returns the full path to the SQLite database file
func (c *Config) DBPath() string {
	return filepath.Join(c.DataDir, "tigclaw.db")
}
