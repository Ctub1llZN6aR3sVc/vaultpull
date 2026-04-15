package env

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Backup creates a timestamped backup of the given .env file.
// Returns the backup path, or an empty string if the file does not exist.
func Backup(envPath string) (string, error) {
	_, err := os.Stat(envPath)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", envPath, err)
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", envPath, err)
	}

	timestamp := time.Now().UTC().Format("20060102T150405Z")
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	backupPath := filepath.Join(dir, fmt.Sprintf("%s.%s.bak", base, timestamp))

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", fmt.Errorf("write backup %s: %w", backupPath, err)
	}

	return backupPath, nil
}

// PruneBackups removes old backups for the given env file, keeping only the
// most recent `keep` backups.
func PruneBackups(envPath string, keep int) error {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	pattern := filepath.Join(dir, fmt.Sprintf("%s.*.bak", base))

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob backups: %w", err)
	}

	if len(matches) <= keep {
		return nil
	}

	// matches from Glob are lexicographically sorted; oldest timestamps first
	toDelete := matches[:len(matches)-keep]
	for _, f := range toDelete {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove backup %s: %w", f, err)
		}
	}
	return nil
}
