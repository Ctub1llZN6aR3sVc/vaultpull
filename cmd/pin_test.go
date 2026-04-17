package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestPinCmd_IsRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "pin" {
			return
		}
	}
	t.Fatal("pin command not registered")
}

func TestPinAdd_CreatesEntry(t *testing.T) {
	dir := t.TempDir()
	pf := filepath.Join(dir, "pins.json")

	buf := &bytes.Buffer{}
	pinAddCmd.SetOut(buf)
	pinFile = pf
	pinTTL = 0

	if err := pinAddCmd.RunE(pinAddCmd, []string{"MY_KEY", "my-value"}); err != nil {
		t.Fatalf("RunE: %v", err)
	}
	if !strings.Contains(buf.String(), "MY_KEY") {
		t.Fatalf("expected output to mention key, got: %s", buf.String())
	}
}

func TestPinList_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	pf := filepath.Join(dir, "pins.json")

	buf := &bytes.Buffer{}
	pinListCmd.SetOut(buf)
	pinFile = pf

	if err := pinListCmd.RunE(pinListCmd, []string{}); err != nil {
		t.Fatalf("RunE: %v", err)
	}
	if !strings.Contains(buf.String(), "No pinned") {
		t.Fatalf("expected empty message, got: %s", buf.String())
	}
}

func TestPinList_ShowsEntries(t *testing.T) {
	dir := t.TempDir()
	pf := filepath.Join(dir, "pins.json")
	pinFile = pf
	pinTTL = 0

	if err := pinAddCmd.RunE(pinAddCmd, []string{"SECRET_KEY", "val"}); err != nil {
		t.Fatalf("add: %v", err)
	}

	buf := &bytes.Buffer{}
	pinListCmd.SetOut(buf)
	pinFile = pf
	if err := pinListCmd.RunE(pinListCmd, []string{}); err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(buf.String(), "SECRET_KEY") {
		t.Fatalf("expected SECRET_KEY in output, got: %s", buf.String())
	}
}
