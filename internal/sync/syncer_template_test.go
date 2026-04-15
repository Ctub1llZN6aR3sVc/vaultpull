package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eliziario/vaultpull/internal/vault"
)

func TestRun_RendersTemplatesInValues(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"BASE_URL": "https://example.com",
			"API_URL":  "${BASE_URL}/v1",
		},
	})

	s := New(client, envFile, []string{"secret/app"}, Options{
		RenderTemplates: true,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	data, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("failed to read env file: %v", err)
	}
	content := string(data)
	if !contains(content, "API_URL=https://example.com/v1") {
		t.Errorf("expected rendered API_URL in output, got:\n%s", content)
	}
}

func TestRun_SkipsTemplateRenderingWhenDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"BASE_URL": "https://example.com",
			"API_URL":  "${BASE_URL}/v1",
		},
	})

	s := New(client, envFile, []string{"secret/app"}, Options{
		RenderTemplates: false,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	data, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("failed to read env file: %v", err)
	}
	content := string(data)
	if !contains(content, "${BASE_URL}/v1") {
		t.Errorf("expected raw template value in output, got:\n%s", content)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
