package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

// MariaDBBackup handles MariaDB database backups
type MariaDBBackup struct {
	config DatabaseConfig
}

// NewMariaDBBackup creates a new MariaDB backup handler
func NewMariaDBBackup(config DatabaseConfig) *MariaDBBackup {
	return &MariaDBBackup{
		config: config,
	}
}

// Backup performs a MariaDB database backup
func (m *MariaDBBackup) Backup(backupDir string) error {
	// Create timestamp for the backup file
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("%s-%s.sql", m.config.Database, timestamp))

	// Construct mysqldump command
	cmd := exec.Command("mysqldump",
		"--host="+m.config.Host,
		fmt.Sprintf("--port=%d", m.config.Port),
		"--user="+m.config.Username,
		"--password="+m.config.Password,
		"--databases",
		m.config.Database,
		"--single-transaction", // Ensures consistent backup
		"--quick",              // Better for large tables
		"--add-drop-database",  // Adds DROP DATABASE statement
		"--add-drop-table",     // Adds DROP TABLE statements
		"--result-file="+backupFile,
	)

	// Execute the backup
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mysqldump failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}
