package env

import (
	"testing"
)

func TestSanitize_NoOptions(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	res := Sanitize(secrets, SanitizeOptions{})
	if res.Sanitized["KEY"] != "value" {
		t.Errorf("expected value unchanged")
	}
	if len(res.ChangedKeys) != 0 {
		t.Errorf("expected no changed keys")
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	secrets := map[string]string{"KEY": "val\x01ue\x00"}
	res := Sanitize(secrets, SanitizeOptions{StripControlChars: true})
	if res.Sanitized["KEY"] != "value" {
		t.Errorf("expected control chars stripped, got %q", res.Sanitized["KEY"])
	}
	if len(res.ChangedKeys) != 1 || res.ChangedKeys[0] != "KEY" {
		t.Errorf("expected KEY in changed keys")
	}
}

func TestSanitize_TrimQuotes(t *testing.T) {
	secrets := map[string]string{
		"A": `"hello"`,
		"B": `'world'`,
		"C": "plain",
	}
	res := Sanitize(secrets, SanitizeOptions{TrimQuotes: true})
	if res.Sanitized["A"] != "hello" {
		t.Errorf("expected double quotes trimmed")
	}
	if res.Sanitized["B"] != "world" {
		t.Errorf("expected single quotes trimmed")
	}
	if res.Sanitized["C"] != "plain" {
		t.Errorf("expected plain unchanged")
	}
}

func TestSanitize_NormalizeNewlines(t *testing.T) {
	secrets := map[string]string{"KEY": "line1\r\nline2\rline3"}
	res := Sanitize(secrets, SanitizeOptions{NormalizeNewlines: true})
	expected := "line1\nline2\nline3"
	if res.Sanitized["KEY"] != expected {
		t.Errorf("expected %q got %q", expected, res.Sanitized["KEY"])
	}
}

func TestSanitize_CombinedOptions(t *testing.T) {
	secrets := map[string]string{"KEY": "\"val\x01ue\""}
	res := Sanitize(secrets, SanitizeOptions{
		StripControlChars: true,
		TrimQuotes:        true,
	})
	if res.Sanitized["KEY"] != "value" {
		t.Errorf("expected combined sanitize, got %q", res.Sanitized["KEY"])
	}
}

func TestSanitize_SummaryNoChanges(t *testing.T) {
	res := SanitizeResult{Sanitized: map[string]string{}, ChangedKeys: nil}
	if res.Summary() != "sanitize: no changes" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestSanitize_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"KEY": `"quoted"`}
	Sanitize(secrets, SanitizeOptions{TrimQuotes: true})
	if secrets["KEY"] != `"quoted"` {
		t.Errorf("input was mutated")
	}
}
