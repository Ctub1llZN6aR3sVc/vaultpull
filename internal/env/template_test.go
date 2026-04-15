package env

import (
	"os"
	"testing"
)

func TestRenderTemplate_NoVariables(t *testing.T) {
	result, err := RenderTemplate("hello world", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("expected 'hello world', got %q", result)
	}
}

func TestRenderTemplate_BraceStyle(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	result, err := RenderTemplate("host=${DB_HOST}", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "host=localhost" {
		t.Errorf("expected 'host=localhost', got %q", result)
	}
}

func TestRenderTemplate_BareStyle(t *testing.T) {
	secrets := map[string]string{"PORT": "5432"}
	result, err := RenderTemplate("port=$PORT", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "port=5432" {
		t.Errorf("expected 'port=5432', got %q", result)
	}
}

func TestRenderTemplate_FallsBackToEnv(t *testing.T) {
	os.Setenv("MY_TEST_VAR", "from-env")
	defer os.Unsetenv("MY_TEST_VAR")

	result, err := RenderTemplate("val=${MY_TEST_VAR}", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "val=from-env" {
		t.Errorf("expected 'val=from-env', got %q", result)
	}
}

func TestRenderTemplate_MissingVariableReturnsError(t *testing.T) {
	_, err := RenderTemplate("${MISSING_KEY}", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestRenderMap_InterpolatesValues(t *testing.T) {
	secrets := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	out, err := RenderMap(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_URL"] != "https://example.com/api" {
		t.Errorf("expected 'https://example.com/api', got %q", out["API_URL"])
	}
}

func TestRenderMap_ReturnsErrorOnMissing(t *testing.T) {
	secrets := map[string]string{
		"API_URL": "${UNDEFINED}/api",
	}
	_, err := RenderMap(secrets)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
