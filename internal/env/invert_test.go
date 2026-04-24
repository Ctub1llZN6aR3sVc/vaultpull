package env

import (
	"testing"
)

func TestInvert_SwapsKeysAndValues(t *testing.T) {
	secrets := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	res, err := Invert(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Inverted["localhost"] != "HOST" {
		t.Errorf("expected localhost->HOST, got %q", res.Inverted["localhost"])
	}
	if res.Inverted["5432"] != "PORT" {
		t.Errorf("expected 5432->PORT, got %q", res.Inverted["5432"])
	}
}

func TestInvert_SkipsEmptyValues(t *testing.T) {
	secrets := map[string]string{
		"EMPTY": "",
		"KEY":   "value",
	}
	res, err := Invert(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
	if _, ok := res.Inverted[""); ok {
		t.Error("empty value should not appear as a key")
	}
}

func TestInvert_ConflictRecorded(t *testing.T) {
	secrets := map[string]string{
		"A": "same",
		"B": "same",
	}
	res, err := Invert(secrets, &InvertOptions{FailOnConflict: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) == 0 {
		t.Error("expected at least one conflict")
	}
}

func TestInvert_FailOnConflict(t *testing.T) {
	secrets := map[string]string{
		"A": "dup",
		"B": "dup",
	}
	_, err := Invert(secrets, &InvertOptions{FailOnConflict: true})
	if err == nil {
		t.Error("expected error on duplicate value with FailOnConflict")
	}
}

func TestInvert_NilOptions(t *testing.T) {
	secrets := map[string]string{"K": "V"}
	res, err := Invert(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Inverted["V"] != "K" {
		t.Errorf("expected V->K")
	}
}

func TestInvert_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"ORIG": "val"}
	_, _ = Invert(secrets, nil)
	if _, ok := secrets["val"]; ok {
		t.Error("Invert must not mutate the input map")
	}
}

func TestInvert_SummaryNoConflicts(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	res, _ := Invert(secrets, nil)
	s := res.Summary()
	if s == "" {
		t.Error("summary should not be empty")
	}
}
