package env

import (
	"strings"
	"testing"
)

func TestSubstitute_NoOptions(t *testing.T) {
	secrets := map[string]string{"A": "${B}", "B": "hello"}
	out, res, err := Substitute(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// nil opts: pass-through, no substitution
	if out["A"] != "${B}" {
		t.Errorf("expected pass-through, got %q", out["A"])
	}
	if len(res.Substituted) != 0 {
		t.Errorf("expected no substitutions, got %v", res.Substituted)
	}
}

func TestSubstitute_ResolvesReference(t *testing.T) {
	secrets := map[string]string{"BASE": "http", "URL": "${BASE}://example.com"}
	out, res, err := Substitute(secrets, &SubstituteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://example.com" {
		t.Errorf("got %q", out["URL"])
	}
	if len(res.Substituted) != 1 || res.Substituted[0] != "URL" {
		t.Errorf("expected URL substituted, got %v", res.Substituted)
	}
}

func TestSubstitute_StrictErrorOnMissing(t *testing.T) {
	secrets := map[string]string{"KEY": "${MISSING}"}
	_, _, err := Substitute(secrets, &SubstituteOptions{Strict: true})
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
	if !strings.Contains(err.Error(), "MISSING") {
		t.Errorf("error should mention missing key, got: %v", err)
	}
}

func TestSubstitute_FallbackForMissing(t *testing.T) {
	secrets := map[string]string{"KEY": "${NOPE}"}
	out, res, err := Substitute(secrets, &SubstituteOptions{Fallback: "default"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "default" {
		t.Errorf("expected fallback value, got %q", out["KEY"])
	}
	if len(res.Unresolved) == 0 {
		t.Error("expected unresolved entry")
	}
}

func TestSubstitute_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{"A": "${B}", "B": "world"}
	out, res, err := Substitute(secrets, &SubstituteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// dry run returns original map
	if out["A"] != "${B}" {
		t.Errorf("dry run should preserve original, got %q", out["A"])
	}
	// but result still reports what would have changed
	if len(res.Substituted) != 1 {
		t.Errorf("expected 1 substitution reported, got %d", len(res.Substituted))
	}
}

func TestSubstitute_SummaryNoUnresolved(t *testing.T) {
	res := SubstituteResult{Substituted: []string{"A", "B"}}
	if !strings.Contains(res.Summary(), "2 substituted") {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestSubstitute_SummaryWithUnresolved(t *testing.T) {
	res := SubstituteResult{Substituted: []string{"A"}, Unresolved: []string{"X"}}
	s := res.Summary()
	if !strings.Contains(s, "unresolved") || !strings.Contains(s, "X") {
		t.Errorf("unexpected summary: %s", s)
	}
}
