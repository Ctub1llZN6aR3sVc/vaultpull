package env

import (
	"strings"
	"testing"
)

func TestInterpolate_NoOptions(t *testing.T) {
	secrets := map[string]string{"A": "hello", "B": "world"}
	out, res, err := Interpolate(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "hello" || out["B"] != "world" {
		t.Errorf("values should be unchanged, got %v", out)
	}
	if len(res.Resolved) != 0 {
		t.Errorf("expected 0 resolved, got %d", len(res.Resolved))
	}
}

func TestInterpolate_ResolvesReference(t *testing.T) {
	secrets := map[string]string{
		"HOST": "localhost",
		"URL":  "http://${HOST}:8080",
	}
	out, res, err := Interpolate(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://localhost:8080" {
		t.Errorf("expected resolved URL, got %q", out["URL"])
	}
	if len(res.Resolved) != 1 || res.Resolved[0] != "URL" {
		t.Errorf("expected URL in resolved list, got %v", res.Resolved)
	}
}

func TestInterpolate_StrictErrorOnMissing(t *testing.T) {
	secrets := map[string]string{"GREETING": "Hello ${NAME}"}
	_, _, err := Interpolate(secrets, &InterpolateOptions{Strict: true})
	if err == nil {
		t.Fatal("expected error for missing reference")
	}
	if !strings.Contains(err.Error(), "NAME") {
		t.Errorf("error should mention missing key NAME, got: %v", err)
	}
}

func TestInterpolate_FallbackForMissing(t *testing.T) {
	secrets := map[string]string{"MSG": "Hello ${UNKNOWN}"}
	out, res, err := Interpolate(secrets, &InterpolateOptions{Fallback: "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MSG"] != "Hello world" {
		t.Errorf("expected fallback substitution, got %q", out["MSG"])
	}
	if len(res.Missing) != 1 || res.Missing[0] != "UNKNOWN" {
		t.Errorf("expected UNKNOWN in missing list, got %v", res.Missing)
	}
}

func TestInterpolate_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{
		"BASE": "base",
		"FULL": "${BASE}_value",
	}
	out, _, err := Interpolate(secrets, &InterpolateOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FULL"] != "${BASE}_value" {
		t.Errorf("dry run should not modify values, got %q", out["FULL"])
	}
}

func TestInterpolate_SelfReferenceIgnored(t *testing.T) {
	secrets := map[string]string{"X": "${X}_suffix"}
	out, _, err := Interpolate(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "${X}_suffix" {
		t.Errorf("self-reference should be left unchanged, got %q", out["X"])
	}
}

func TestInterpolate_SummaryNoErrors(t *testing.T) {
	res := InterpolateResult{Resolved: []string{"A", "B"}}
	if !strings.Contains(res.Summary(), "2 key(s)") {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
	if !strings.Contains(res.Summary(), "no missing") {
		t.Errorf("expected no missing in summary: %s", res.Summary())
	}
}
