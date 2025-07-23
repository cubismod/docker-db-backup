# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a Go-based database backup tool designed to run in Docker containers. The application:

- Reads database configurations from a YAML file (`config.yaml`)
- Supports PostgreSQL, MariaDB/MySQL, and Redis databases
- Creates timestamped backup files using native database tools (`pg_dump`, `mariadb-dump`, `redis-cli`)
- Uses a plugin-like architecture with separate backup handlers for each database type

### Core Components

- `main.go`: Entry point that loads config, creates backup directory, and orchestrates backups
- `config.go`: Configuration loading and data structures for YAML parsing
- `postgres.go`: PostgreSQL backup implementation using `pg_dump`
- `mariadb.go`: MariaDB/MySQL backup implementation using `mariadb-dump`
- `redis.go`: Redis backup implementation using `redis-cli --rdb`
- `Dockerfile`: Multi-stage build that includes database client tools

## Development Commands

### Building and Running

```bash
# Build the Go binary
go build -o db-backup .

# Run with default config
./db-backup

# Run with custom config file
./db-backup path/to/config.yaml

# Build Docker image
docker build -t docker-db-backup .

# Run in Docker
docker run -v $(pwd)/backups:/app/backups -v $(pwd)/config.yaml:/app/config.yaml docker-db-backup
```

### Go Module Management

```bash
# Download dependencies
go mod download

# Update dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Configuration

The application expects a `config.yaml` file with:
- `backup_dir`: Directory where backup files are stored
- `databases`: Array of database configurations with type, connection details, and credentials

Database types supported: `postgres`, `mariadb`, `mysql`, `redis`

## External Dependencies

The application requires database client tools to be available:
- PostgreSQL: `pg_dump`
- MariaDB/MySQL: `mariadb-dump` (at `/usr/bin/mariadb-dump`)
- Redis: `redis-cli`

These are installed in the Docker image via Alpine packages (`postgresql-client`, `mariadb-client`, `redis`).