package env

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions configures the rotation behaviour.
type RotateOptions struct {
	// BackupDir is the directory where rotated files are stored.
	// Defaults to the same directory as the source file.
	BackupDir string
	// MaxBackups is the maximum number of rotated copies to keep (0 = unlimited).
	MaxBackups int
}

// RotateResult holds metadata about a completed rotation.
type RotateResult struct {
	BackupPath string
	RotatedAt  time.Time
}

// Rotate backs up the current env file, writes the new secrets map to it,
// and prunes old backups according to opts.MaxBackups.
// If the file does not yet exist it is created without a backup step.
func Rotate(path string, secrets map[string]string, opts RotateOptions) (RotateResult, error) {
	result := RotateResult{RotatedAt: time.Now()}

	backupDir := opts.BackupDir
	if backupDir == "" {
		backupDir = filepath.Dir(path)
	}

	if _, err := os.Stat(path); err == nil {
		ts := result.RotatedAt.UTC().Format("20060102T150405Z")
		base := filepath.Base(path)
		backupPath := filepath.Join(backupDir, fmt.Sprintf("%s.%s.bak", base, ts))

		if err := Backup(path, backupPath); err != nil {
			return result, fmt.Errorf("rotate: backup failed: %w", err)
		}
		result.BackupPath = backupPath

		if opts.MaxBackups > 0 {
			if err := PruneBackups(backupDir, base, opts.MaxBackups); err != nil {
				return result, fmt.Errorf("rotate: prune failed: %w", err)
			}
		}
	}

	w := NewWriter(path)
	if err := w.Write(secrets); err != nil {
		return result, fmt.Errorf("rotate: write failed: %w", err)
	}

	return result, nil
}
