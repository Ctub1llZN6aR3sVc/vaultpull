package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintDiff_NoChanges(t *testing.T) {
	result := &DiffResult{}
	var buf bytes.Buffer
	PrintDiff(&buf, result)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got: %s", buf.String())
	}
}

func TestPrintDiff_ShowsAdded(t *testing.T) {
	result := &DiffResult{
		Changes: []Change{
			{Key: "NEW_KEY", Type: Added, NewValue: "val"},
		},
	}
	var buf bytes.Buffer
	PrintDiff(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected '+ NEW_KEY', got: %s", out)
	}
	if !strings.Contains(out, "+1 added") {
		t.Errorf("expected summary with +1 added, got: %s", out)
	}
}

func TestPrintDiff_ShowsRemoved(t *testing.T) {
	result := &DiffResult{
		Changes: []Change{
			{Key: "OLD_KEY", Type: Removed, OldValue: "val"},
		},
	}
	var buf bytes.Buffer
	PrintDiff(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected '- OLD_KEY', got: %s", out)
	}
}

func TestPrintDiff_ShowsChanged(t *testing.T) {
	result := &DiffResult{
		Changes: []Change{
			{Key: "EXISTING", Type: Changed, OldValue: "old", NewValue: "new"},
		},
	}
	var buf bytes.Buffer
	PrintDiff(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "~ EXISTING") {
		t.Errorf("expected '~ EXISTING', got: %s", out)
	}
}

func TestPrintDiff_SortedOutput(t *testing.T) {
	result := &DiffResult{
		Changes: []Change{
			{Key: "Z_KEY", Type: Added, NewValue: "z"},
			{Key: "A_KEY", Type: Added, NewValue: "a"},
		},
	}
	var buf bytes.Buffer
	PrintDiff(&buf, result)
	out := buf.String()
	aIdx := strings.Index(out, "A_KEY")
	zIdx := strings.Index(out, "Z_KEY")
	if aIdx > zIdx {
		t.Errorf("expected A_KEY before Z_KEY in output")
	}
}
