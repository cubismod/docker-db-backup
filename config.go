package main

import (
	"fmt"
	"os"

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
	Databases []DatabaseConfig `yaml:"databases"`
	BackupDir string           `yaml:"backup_dir"`
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

	return &config, nil
}

// GetFilenamePrefix returns the prefix to use for backup filenames
// Uses alias if provided, otherwise falls back to database name
func (d *DatabaseConfig) GetFilenamePrefix() string {
	if d.Alias != "" {
		return d.Alias
	}
	return d.Database
}
