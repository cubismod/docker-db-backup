package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// RedisBackup handles Redis database backups
type RedisBackup struct {
	config DatabaseConfig
}

// NewRedisBackup creates a new Redis backup handler
func NewRedisBackup(config DatabaseConfig) *RedisBackup {
	return &RedisBackup{
		config: config,
	}
}

// Backup performs a Redis database backup using redis-cli --rdb
func (r *RedisBackup) Backup(backupDir string) error {
	// Create timestamp for the backup file
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("%s-%s.rdb", r.config.Database, timestamp))

	// Build redis-cli command arguments
	args := []string{
		"-h", r.config.Host,
		"-p", strconv.Itoa(r.config.Port),
	}

	// Add authentication if password is provided
	if r.config.Password != "" {
		args = append(args, "-a", r.config.Password)
	}

	// Add database selection if specified (Redis databases are numbered 0-15)
	if r.config.Database != "" {
		args = append(args, "-n", r.config.Database)
	}

	// Add the --rdb flag with output file
	args = append(args, "--rdb", backupFile)

	cmd := exec.Command("redis-cli", args...)

	// Execute the backup
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("redis-cli --rdb failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}