package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditEntry records a single sync event.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Profile   string    `json:"profile"`
	EnvFile   string    `json:"env_file"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
	Changed   int       `json:"changed"`
}

// WriteAuditLog appends an audit entry as a JSON line to the given log file.
// The parent directory is created if it does not exist.
func WriteAuditLog(logPath string, entry AuditEntry) error {
	if err := os.MkdirAll(filepath.Dir(logPath), 0o700); err != nil {
		return fmt.Errorf("audit: create log dir: %w", err)
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}

// ReadAuditLog reads all audit entries from the given log file.
// Returns an empty slice if the file does not exist.
func ReadAuditLog(logPath string) ([]AuditEntry, error) {
	f, err := os.Open(logPath)
	if os.IsNotExist(err) {
		return []AuditEntry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	var entries []AuditEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e AuditEntry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
