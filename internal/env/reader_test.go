package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	return path
}

func TestRead_ParsesKeyValuePairs(t *testing.T) {
	path := writeEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\n")

	got, err := Read(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", got["DB_PORT"])
	}
}

func TestRead_IgnoresComments(t *testing.T) {
	path := writeEnvFile(t, "# this is a comment\nFOO=bar\n")

	got, err := Read(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Errorf("expected 1 key, got %d", len(got))
	}
	if got["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", got["FOO"])
	}
}

func TestRead_UnquotesValues(t *testing.T) {
	path := writeEnvFile(t, `SECRET="hello world"
TOKEN='abc123'
`)

	got, err := Read(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["SECRET"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", got["SECRET"])
	}
	if got["TOKEN"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", got["TOKEN"])
	}
}

func TestRead_MissingFileReturnsEmpty(t *testing.T) {
	got, err := Read("/nonexistent/path/.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestUnquoteValue(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{`"quoted"`, "quoted"},
		{`'single'`, "single"},
		{`plain`, "plain"},
		{`"`, `"`},
	}

	for _, tc := range cases {
		got := unquoteValue(tc.input)
		if got != tc.expected {
			t.Errorf("unquoteValue(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
