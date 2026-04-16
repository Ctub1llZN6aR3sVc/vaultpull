package env

import "testing"

func TestRedact_ExplicitKeys(t *testing.T) {
	secrets := map[string]string{
		"API_KEY": "abc123",
		"HOST":    "localhost",
	}
	out := Redact(secrets, RedactOptions{Keys: []string{"API_KEY"}})
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %s", out["API_KEY"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", out["HOST"])
	}
}

func TestRedact_AutoDetect(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "secret",
		"APP_NAME":   "myapp",
	}
	out := Redact(secrets, RedactOptions{AutoDetect: true})
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %s", out["DB_PASSWORD"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected myapp, got %s", out["APP_NAME"])
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	secrets := map[string]string{"TOKEN": "tok_live_abc"}
	out := Redact(secrets, RedactOptions{Keys: []string{"TOKEN"}, Placeholder: "***"})
	if out["TOKEN"] != "***" {
		t.Errorf("expected ***, got %s", out["TOKEN"])
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"SECRET_KEY": "original"}
	Redact(secrets, RedactOptions{AutoDetect: true})
	if secrets["SECRET_KEY"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestRedact_EmptySecrets(t *testing.T) {
	out := Redact(map[string]string{}, RedactOptions{AutoDetect: true})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}

func TestRedact_CaseInsensitiveKeys(t *testing.T) {
	secrets := map[string]string{"api_key": "value"}
	out := Redact(secrets, RedactOptions{Keys: []string{"API_KEY"}})
	if out["api_key"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %s", out["api_key"])
	}
}
