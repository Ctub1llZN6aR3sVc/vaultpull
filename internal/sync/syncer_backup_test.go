package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_CreatesBackupWhenEnabled(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	// Pre-populate the env file so there is something to back up.
	require.NoError(t, os.WriteFile(envPath, []byte("EXISTING=yes\n"), 0600))

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"TOKEN": "abc123"},
	})

	syncer := New(client, envPath, []string{"secret/app"}, true)
	_, err := syncer.Run()
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(dir, ".env.*.bak"))
	require.NoError(t, err)
	assert.Len(t, matches, 1, "expected exactly one backup file")

	data, err := os.ReadFile(matches[0])
	require.NoError(t, err)
	assert.Contains(t, string(data), "EXISTING=yes")
}

func TestRun_NoBackupWhenDisabled(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("OLD=value\n"), 0600))

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"NEW": "val"},
	})

	syncer := New(client, envPath, []string{"secret/app"}, false)
	_, err := syncer.Run()
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(dir, ".env.*.bak"))
	require.NoError(t, err)
	assert.Empty(t, matches, "expected no backup files")
}

func TestRun_BackupAbsentFileIsNoop(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env") // does not exist yet

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"KEY": "value"},
	})

	syncer := New(client, envPath, []string{"secret/app"}, true)
	_, err := syncer.Run()
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(dir, ".env.*.bak"))
	require.NoError(t, err)
	assert.Empty(t, matches, "no backup expected when original file was absent")
}
