package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// PostgresBackup handles PostgreSQL database backups
type PostgresBackup struct {
	config DatabaseConfig
}

// NewPostgresBackup creates a new PostgreSQL backup handler
func NewPostgresBackup(config DatabaseConfig) *PostgresBackup {
	return &PostgresBackup{
		config: config,
	}
}

// Backup performs a PostgreSQL database backup
func (p *PostgresBackup) Backup(backupDir string) error {
	// Create timestamp for the backup file
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("%s-%s.sql", p.config.GetFilenamePrefix(), timestamp))

	cmd := exec.Command("pg_dump",
		"-h", p.config.Host,
		"-p", fmt.Sprintf("%d", p.config.Port),
		"-U", p.config.Username,
		"-d", p.config.Database,
		"-F", "c", // Custom format (compressed)
		"-f", backupFile,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.config.Password))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}
