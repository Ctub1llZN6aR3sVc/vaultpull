package env

import (
	"strings"
	"testing"
)

func TestLint_CleanSecrets(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"API_KEY": "abc123",
	}
	result := Lint(secrets)
	if !result.IsClean() {
		t.Fatalf("expected clean, got: %v", result.Issues)
	}
}

func TestLint_KeyWithLeadingWhitespace(t *testing.T) {
	secrets := map[string]string{" DB_HOST": "localhost"}
	result := Lint(secrets)
	if result.IsClean() {
		t.Fatal("expected lint issue for leading whitespace")
	}
	if !strings.Contains(result.Issues[0].Message, "whitespace") {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_KeyWithSpaces(t *testing.T) {
	secrets := map[string]string{"DB HOST": "localhost"}
	result := Lint(secrets)
	found := false
	for _, issue := range result.Issues {
		if strings.Contains(issue.Message, "spaces") {
			found = true
		}
	}
	if !found {
		t.Error("expected lint issue for spaces in key")
	}
}

func TestLint_OversizedValue(t *testing.T) {
	secrets := map[string]string{"BIG": strings.Repeat("x", 5000)}
	result := Lint(secrets)
	if result.IsClean() {
		t.Fatal("expected lint issue for oversized value")
	}
	if !strings.Contains(result.Issues[0].Message, "4096") {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_SummaryWithIssues(t *testing.T) {
	secrets := map[string]string{" BAD KEY": "val"}
	result := Lint(secrets)
	if result.IsClean() {
		t.Fatal("expected issues")
	}
	if !strings.Contains(result.Summary(), "issue") {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestLint_SummaryClean(t *testing.T) {
	result := Lint(map[string]string{"OK": "val"})
	if result.Summary() != "no lint issues found" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}
