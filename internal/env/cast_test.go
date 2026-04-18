package env

import (
	"testing"
)

func TestCast_NoOptions(t *testing.T) {
	secrets := map[string]string{"PORT": "8080", "NAME": "app"}
	out, results, err := Cast(secrets, CastOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080")
	}
}

func TestCast_IntType(t *testing.T) {
	secrets := map[string]string{"PORT": "8080.0"}
	out, results, err := Cast(secrets, CastOptions{Types: map[string]string{"PORT": "int"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected 8080, got %s", out["PORT"])
	}
	if len(results) != 1 || results[0].Type != "int" {
		t.Errorf("expected one int result")
	}
}

func TestCast_BoolType(t *testing.T) {
	secrets := map[string]string{"ENABLED": "true", "DEBUG": "1"}
	out, _, err := Cast(secrets, CastOptions{Types: map[string]string{"ENABLED": "bool", "DEBUG": "bool"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENABLED"] != "true" {
		t.Errorf("expected true")
	}
	if out["DEBUG"] != "true" {
		t.Errorf("expected true for 1")
	}
}

func TestCast_FloatType(t *testing.T) {
	secrets := map[string]string{"RATIO": "3.14"}
	out, _, err := Cast(secrets, CastOptions{Types: map[string]string{"RATIO": "float"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["RATIO"] != "3.14" {
		t.Errorf("expected 3.14, got %s", out["RATIO"])
	}
}

func TestCast_StrictErrorOnInvalid(t *testing.T) {
	secrets := map[string]string{"PORT": "notanint"}
	_, _, err := Cast(secrets, CastOptions{
		Types:  map[string]string{"PORT": "int"},
		Strict: true,
	})
	if err == nil {
		t.Fatal("expected error in strict mode")
	}
}

func TestCast_NonStrictSkipsInvalid(t *testing.T) {
	secrets := map[string]string{"PORT": "notanint"}
	out, results, err := Cast(secrets, CastOptions{
		Types:  map[string]string{"PORT": "int"},
		Strict: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for failed cast")
	}
	if out["PORT"] != "notanint" {
		t.Errorf("expected original value preserved")
	}
}

func TestCast_MissingKeyIgnored(t *testing.T) {
	secrets := map[string]string{"NAME": "app"}
	_, results, err := Cast(secrets, CastOptions{Types: map[string]string{"PORT": "int"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for missing key")
	}
}
