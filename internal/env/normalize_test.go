package env

import (
	"testing"
)

func TestNormalize_NoOptions(t *testing.T) {
	secrets := map[string]string{"my-key": "value", "OTHER": "val"}
	res := Normalize(secrets, NormalizeOptions{})
	if res.Output["my-key"] != "value" {
		t.Errorf("expected key unchanged")
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected no renames")
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost", "DB_PORT": "5432"}
	res := Normalize(secrets, NormalizeOptions{UppercaseKeys: true})
	if _, ok := res.Output["DB_HOST"]; !ok {
		t.Errorf("expected DB_HOST in output")
	}
	if _, ok := res.Output["DB_PORT"]; !ok {
		t.Errorf("expected DB_PORT in output")
	}
}

func TestNormalize_ReplaceHyphens(t *testing.T) {
	secrets := map[string]string{"my-secret-key": "val"}
	res := Normalize(secrets, NormalizeOptions{ReplaceHyphens: true})
	if _, ok := res.Output["my_secret_key"]; !ok {
		t.Errorf("expected hyphen replaced with underscore")
	}
	if len(res.Renamed) != 1 {
		t.Errorf("expected 1 rename, got %d", len(res.Renamed))
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	secrets := map[string]string{"KEY": "  spaced  "}
	res := Normalize(secrets, NormalizeOptions{TrimValues: true})
	if res.Output["KEY"] != "spaced" {
		t.Errorf("expected trimmed value, got %q", res.Output["KEY"])
	}
}

func TestNormalize_CombinedOptions(t *testing.T) {
	secrets := map[string]string{"my-key": "  hello  "}
	res := Normalize(secrets, NormalizeOptions{UppercaseKeys: true, ReplaceHyphens: true, TrimValues: true})
	if res.Output["MY_KEY"] != "hello" {
		t.Errorf("expected MY_KEY=hello, got %v", res.Output)
	}
}

func TestNormalize_SummaryNoChanges(t *testing.T) {
	res := NormalizeResult{Output: map[string]string{"A": "1"}}
	if res.Summary() != "normalize: no keys altered" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestNormalize_SummaryWithChanges(t *testing.T) {
	res := NormalizeResult{
		Output:  map[string]string{"A_B": "1"},
		Renamed: []string{"a-b→A_B"},
	}
	if res.Summary() == "normalize: no keys altered" {
		t.Errorf("expected non-empty summary")
	}
}
