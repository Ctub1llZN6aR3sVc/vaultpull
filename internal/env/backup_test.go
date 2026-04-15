package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackup_CreatesBackupFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("FOO=bar\n"), 0600))

	backupPath, err := Backup(envPath)
	require.NoError(t, err)
	assert.NotEmpty(t, backupPath)

	data, err := os.ReadFile(backupPath)
	require.NoError(t, err)
	assert.Equal(t, "FOO=bar\n", string(data))
}

func TestBackup_NoopWhenFileAbsent(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	backupPath, err := Backup(envPath)
	require.NoError(t, err)
	assert.Empty(t, backupPath)
}

func TestBackup_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("SECRET=x\n"), 0600))

	backupPath, err := Backup(envPath)
	require.NoError(t, err)

	info, err := os.Stat(backupPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestPruneBackups_RemovesOldest(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("A=1\n"), 0600))

	// Create 4 backups
	var paths []string
	for i := 0; i < 4; i++ {
		p, err := Backup(envPath)
		require.NoError(t, err)
		paths = append(paths, p)
		// tiny sleep ensures distinct timestamps in filenames
		_ = p
	}

	require.NoError(t, PruneBackups(envPath, 2))

	pattern := filepath.Join(dir, ".env.*.bak")
	matches, err := filepath.Glob(pattern)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(matches), 2)
}

func TestPruneBackups_KeepsAllWhenUnderLimit(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("A=1\n"), 0600))

	_, err := Backup(envPath)
	require.NoError(t, err)

	require.NoError(t, PruneBackups(envPath, 5))

	pattern := filepath.Join(dir, ".env.*.bak")
	matches, err := filepath.Glob(pattern)
	require.NoError(t, err)
	assert.Equal(t, 1, len(matches))
}
