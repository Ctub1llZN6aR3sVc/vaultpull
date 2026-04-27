package env

import (
	"testing"
)

func TestCompute_NoOptions(t *testing.T) {
	secrets := map[string]string{"A": "hello", "B": "world"}
	out, res, err := Compute(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if len(res.Added) != 0 {
		t.Errorf("expected no added keys")
	}
}

func TestCompute_ConcatKeys(t *testing.T) {
	secrets := map[string]string{"FIRST": "hello", "SECOND": "world"}
	opts := &ComputeOptions{
		Expressions: map[string]string{"FULL": "FIRST + SECOND"},
	}
	out, res, err := Compute(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FULL"] != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", out["FULL"])
	}
	if len(res.Added) != 1 || res.Added[0] != "FULL" {
		t.Errorf("expected FULL in added, got %v", res.Added)
	}
}

func TestCompute_NumericSubtract(t *testing.T) {
	secrets := map[string]string{"X": "10", "Y": "3"}
	opts := &ComputeOptions{
		Expressions: map[string]string{"DIFF": "X - Y"},
	}
	out, _, err := Compute(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DIFF"] != "7" {
		t.Errorf("expected '7', got %q", out["DIFF"])
	}
}

func TestCompute_NumericMultiply(t *testing.T) {
	secrets := map[string]string{"A": "6", "B": "7"}
	opts := &ComputeOptions{
		Expressions: map[string]string{"PRODUCT": "A * B"},
	}
	out, _, err := Compute(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PRODUCT"] != "42" {
		t.Errorf("expected '42', got %q", out["PRODUCT"])
	}
}

func TestCompute_LenExpression(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abcdef"}
	opts := &ComputeOptions{
		Expressions: map[string]string{"TOKEN_LEN": "len(TOKEN)"},
	}
	out, _, err := Compute(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN_LEN"] != "6" {
		t.Errorf("expected '6', got %q", out["TOKEN_LEN"])
	}
}

func TestCompute_SkipsExistingWithoutOverwrite(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "existing"}
	opts := &ComputeOptions{
		Expressions: map[string]string{"C": "A + B"},
		Overwrite:   false,
	}
	out, res, err := Compute(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["C"] != "existing" {
		t.Errorf("expected 'existing', got %q", out["C"])
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %v", res.Skipped)
	}
}

func TestCompute_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{"A": "hello", "B": "world"}
	opts := &ComputeOptions{
		Expressions: map[string]string{"C": "A + B"},
		DryRun:      true,
	}
	out, res, err := Compute(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := out["C"]; exists {
		t.Errorf("dry run should not write key C")
	}
	if len(res.Added) != 1 {
		t.Errorf("expected C in added even on dry run")
	}
}

func TestCompute_MissingKeyReturnsError(t *testing.T) {
	secrets := map[string]string{"A": "hello"}
	opts := &ComputeOptions{
		Expressions: map[string]string{"C": "A + MISSING"},
	}
	_, res, err := Compute(secrets, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if len(res.Errors) == 0 {
		t.Errorf("expected errors in result")
	}
}

func TestCompute_SummaryNoErrors(t *testing.T) {
	res := ComputeResult{Added: []string{"X"}, Skipped: []string{}}
	s := res.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
