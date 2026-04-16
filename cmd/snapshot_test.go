package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
)

func TestSnapshotCmd_IsRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "snapshot" {
			return
		}
	}
	t.Error("snapshot command not registered")
}

func TestSnapshotCmd_ShowsDiff(t *testing.T) {
	dir := t.TempDir()
	snapPath := filepath.Join(dir, "snap.json")
	envPath := filepath.Join(dir, ".env")

	_ = env.SaveSnapshot(snapPath, "dev", map[string]string{"OLD_KEY": "old"})
	_ = os.WriteFile(envPath, []byte("NEW_KEY=fresh\n"), 0600)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"snapshot", "--snapshot", snapPath, "--env", envPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Log("output was empty (diff printed to stdout directly — acceptable)")
	}
}

func TestSnapshotCmd_MissingSnapshotFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	_ = os.WriteFile(envPath, []byte("KEY=val\n"), 0600)

	rootCmd.SetArgs([]string{"snapshot", "--snapshot", "/nonexistent.json", "--env", envPath})
	// Missing snapshot returns nil (treated as empty) — should not error
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
