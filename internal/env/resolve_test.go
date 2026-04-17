package env

import (
	"os"
	"testing"
)

func TestResolve_AppliesDefaults(t *testing.T) {
	result, err := Resolve(map[string]string{}, ResolveOptions{
		Defaults: map[string]string{"FOO": "bar"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Secrets["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", result.Secrets["FOO"])
	}
}

func TestResolve_SecretsOverrideDefaults(t *testing.T) {
	result, err := Resolve(map[string]string{"FOO": "vault"}, ResolveOptions{
		Defaults: map[string]string{"FOO": "default"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Secrets["FOO"] != "vault" {
		t.Errorf("expected FOO=vault, got %q", result.Secrets["FOO"])
	}
}

func TestResolve_MissingRequiredReturnsError(t *testing.T) {
	_, err := Resolve(map[string]string{}, ResolveOptions{
		Required: []string{"MUST_EXIST"},
	})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestResolve_FallbackToEnv(t *testing.T) {
	t.Setenv("MY_SECRET", "from-env")
	result, err := Resolve(map[string]string{}, ResolveOptions{
		FallbackToEnv: true,
		Required:      []string{"MY_SECRET"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Secrets["MY_SECRET"] != "from-env" {
		t.Errorf("expected from-env, got %q", result.Secrets["MY_SECRET"])
	}
	if len(result.Warnings) == 0 {
		t.Error("expected warning about OS env fallback")
	}
}

func TestResolve_FallbackToEnvMissingStillErrors(t *testing.T) {
	os.Unsetenv("ABSENT_KEY")
	_, err := Resolve(map[string]string{}, ResolveOptions{
		FallbackToEnv: true,
		Required:      []string{"ABSENT_KEY"},
	})
	if err == nil {
		t.Fatal("expected error when env var also absent")
	}
}

func TestResolve_NoRequiredNoError(t *testing.T) {
	result, err := Resolve(map[string]string{"A": "1"}, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Secrets["A"] != "1" {
		t.Errorf("expected A=1")
	}
}
