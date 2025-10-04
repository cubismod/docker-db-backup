package main

import (
	"fmt"
	"os"
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

// Backup performs a Redis database backup using redis-cli with RESP format
func (r *RedisBackup) Backup(backupDir string) error {
	// Create timestamp for the backup file
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("%s-%s.resp", r.config.GetFilenamePrefix(), timestamp))

	// Build redis-cli command arguments for RESP export
	args := []string{
		"-h", r.config.Host,
		"-p", strconv.Itoa(r.config.Port),
		"--raw",
	}

	// Add authentication if password is provided
	if r.config.Password != "" {
		args = append(args, "-a", r.config.Password)
	}

	// Add database selection if specified (Redis databases are numbered 0-15)
	if r.config.Database != "" {
		args = append(args, "-n", r.config.Database)
	}

	// Use a Lua script to export all keys and values as RESP commands
	luaScript := `
local keys = redis.call('keys', '*')
local result = {}
for i=1,#keys do
    local key = keys[i]
    local keytype = redis.call('type', key)['ok']
    local ttl = redis.call('ttl', key)
    
    if keytype == 'string' then
        local value = redis.call('get', key)
        table.insert(result, '*3\r\n$3\r\nSET\r\n$' .. #key .. '\r\n' .. key .. '\r\n$' .. #value .. '\r\n' .. value .. '\r\n')
    elseif keytype == 'hash' then
        local hash = redis.call('hgetall', key)
        local cmd = '*' .. (#hash + 2) .. '\r\n$5\r\nHMSET\r\n$' .. #key .. '\r\n' .. key .. '\r\n'
        for j=1,#hash,2 do
            cmd = cmd .. '$' .. #hash[j] .. '\r\n' .. hash[j] .. '\r\n$' .. #hash[j+1] .. '\r\n' .. hash[j+1] .. '\r\n'
        end
        table.insert(result, cmd)
    elseif keytype == 'list' then
        local list = redis.call('lrange', key, 0, -1)
        for j=1,#list do
            table.insert(result, '*3\r\n$5\r\nLPUSH\r\n$' .. #key .. '\r\n' .. key .. '\r\n$' .. #list[j] .. '\r\n' .. list[j] .. '\r\n')
        end
    elseif keytype == 'set' then
        local set = redis.call('smembers', key)
        for j=1,#set do
            table.insert(result, '*3\r\n$4\r\nSADD\r\n$' .. #key .. '\r\n' .. key .. '\r\n$' .. #set[j] .. '\r\n' .. set[j] .. '\r\n')
        end
    elseif keytype == 'zset' then
        local zset = redis.call('zrange', key, 0, -1, 'withscores')
        for j=1,#zset,2 do
            table.insert(result, '*4\r\n$4\r\nZADD\r\n$' .. #key .. '\r\n' .. key .. '\r\n$' .. #zset[j+1] .. '\r\n' .. zset[j+1] .. '\r\n$' .. #zset[j] .. '\r\n' .. zset[j] .. '\r\n')
        end
    end
    
    if ttl > 0 then
        table.insert(result, '*3\r\n$6\r\nEXPIRE\r\n$' .. #key .. '\r\n' .. key .. '\r\n$' .. #tostring(ttl) .. '\r\n' .. tostring(ttl) .. '\r\n')
    end
end
return table.concat(result)
`
	args = append(args, "--eval", luaScript, "0")

	cmd := exec.Command("redis-cli", args...)

	// Execute the command and save output to file
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("redis-cli eval failed: %v", err)
	}

	// Write the output to the backup file
	err = os.WriteFile(backupFile, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}

	return nil
}
