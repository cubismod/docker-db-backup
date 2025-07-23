package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting database backup service...")

	// Load configuration
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
		log.Fatalf("Failed to create backup directory: %v", err)
	}

	log.Printf("Loaded configuration for %d databases", len(config.Databases))
	log.Printf("Backup directory: %s", config.BackupDir)

	// Process each database
	for _, dbConfig := range config.Databases {
		log.Printf("Processing database: %s (%s)", dbConfig.Database, dbConfig.Type)

		switch dbConfig.Type {
		case "postgres":
			backup := NewPostgresBackup(dbConfig)
			if err := backup.Backup(config.BackupDir); err != nil {
				log.Printf("Failed to backup database %s: %v", dbConfig.Database, err)
				continue
			}
			log.Printf("Successfully backed up database: %s", dbConfig.Database)
		case "mariadb":
			backup := NewMariaDBBackup(dbConfig)
			if err := backup.Backup(config.BackupDir); err != nil {
				log.Printf("Failed to backup database %s: %v", dbConfig.Database, err)
				continue
			}
			log.Printf("Successfully backed up database: %s", dbConfig.Database)
		case "redis":
			backup := NewRedisBackup(dbConfig)
			if err := backup.Backup(config.BackupDir); err != nil {
				log.Printf("Failed to backup database %s: %v", dbConfig.Database, err)
				continue
			}
			log.Printf("Successfully backed up database: %s", dbConfig.Database)
		default:
			log.Printf("Unsupported database type: %s", dbConfig.Type)
		}
	}

	log.Println("Backup process completed")
}
