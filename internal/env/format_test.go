package env

import (
	"strings"
	"testing"
)

func TestSerialize_DotenvFormat(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Serialize(secrets, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar\n") {
		t.Errorf("expected FOO=bar, got: %s", out)
	}
	if !strings.Contains(out, "BAZ=qux\n") {
		t.Errorf("expected BAZ=qux, got: %s", out)
	}
}

func TestSerialize_ExportFormat(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abc123"}
	out, err := Serialize(secrets, FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export TOKEN=abc123\n") {
		t.Errorf("expected export TOKEN=abc123, got: %s", out)
	}
}

func TestSerialize_JSONFormat(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	out, err := Serialize(secrets, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"KEY": "value"`) {
		t.Errorf("expected JSON key-value, got: %s", out)
	}
}

func TestSerialize_UnknownFormat(t *testing.T) {
	_, err := Serialize(map[string]string{}, Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestSerialize_QuotesSpecialChars(t *testing.T) {
	secrets := map[string]string{"MSG": "hello world"}
	out, err := Serialize(secrets, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestSerialize_SortedOutput(t *testing.T) {
	secrets := map[string]string{"ZZZ": "1", "AAA": "2", "MMM": "3"}
	out, err := Serialize(secrets, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if !strings.HasPrefix(lines[0], "AAA") {
		t.Errorf("expected AAA first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZZZ") {
		t.Errorf("expected ZZZ last, got: %s", lines[2])
	}
}
