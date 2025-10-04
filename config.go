package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Type      string `yaml:"type"` // postgres, mysql, mariadb, etc.
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
	Alias     string `yaml:"alias"`     // optional alias for backup filename
	EnableSSL string `yaml:"enablessl"` // needs to be 0 or 1
}

// Config represents the main configuration structure
type Config struct {
	Databases          []DatabaseConfig `yaml:"databases"`
	BackupDir          string           `yaml:"backup_dir"`
	PrePingURLs        []string         `yaml:"pre_ping_urls,omitempty"`        // URLs to ping before the backup run
	PostPingURLs       []string         `yaml:"post_ping_urls,omitempty"`       // URLs to ping after the backup run
	PingTimeoutSeconds int              `yaml:"ping_timeout_seconds,omitempty"` // optional timeout in seconds for each ping (default 10)
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Apply a sensible default if not set
	if config.PingTimeoutSeconds <= 0 {
		config.PingTimeoutSeconds = 10
	}

	return &config, nil
}

// PingURLs sends HTTP GET requests to the provided URLs with a timeout specified in the config.
// It returns an error if any of the pings fail. The error will list the failing URLs and statuses.
func (c *Config) PingURLs(urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	timeout := time.Duration(c.PingTimeoutSeconds) * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var failures []string
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		if err := pingURL(ctx, client, u); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", u, err))
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf("ping failures: %s", strings.Join(failures, "; "))
	}
	return nil
}

// pingURL performs a single HTTP GET to the given URL using the provided context and client.
// Success is defined as receiving any 2xx or 3xx status code. Non-2xx/3xx responses are considered failures.
func pingURL(ctx context.Context, client *http.Client, url string) error {
	// Create a request bound to the provided context so overall timeout is respected.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}

	return fmt.Errorf("unexpected status %d", resp.StatusCode)
}

// RunPrePings pings all configured pre-run URLs. It returns an error if any ping fails.
func (c *Config) RunPrePings() error {
	return c.PingURLs(c.PrePingURLs)
}

// RunPostPings pings all configured post-run URLs. It returns an error if any ping fails.
func (c *Config) RunPostPings() error {
	return c.PingURLs(c.PostPingURLs)
}

// GetFilenamePrefix returns the prefix to use for backup filenames
// Uses alias if provided, otherwise falls back to database name
func (d *DatabaseConfig) GetFilenamePrefix() string {
	if d.Alias != "" {
		return d.Alias
	}
	return d.Database
}
