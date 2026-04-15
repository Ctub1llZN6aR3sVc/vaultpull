package env

import "fmt"

// ChangeType represents the type of change detected.
type ChangeType string
(
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// DiffResult holds the full set of changes between two env maps.
type DiffResult struct {
	Changes []Change
}

// IsEmpty returns true if there are no changes.
func (d *DiffResult) IsEmpty() bool {
	return len(d.Changes) == 0
}

// Summary returns a human-readable summary of changes.
func (d *DiffResult) Summary() string {
	if d.IsEmpty() {
		return "no changes"
	}
	var out string
	for _, c := range d.Changes {
		switch c.Type {
		case Added:
			out += fmt.Sprintf("+ %s=%q\n", c.Key, c.NewValue)
		case Removed:
			out += fmt.Sprintf("- %s=%q\n", c.Key, c.OldValue)
		case Changed:
			out += fmt.Sprintf("~ %s: %q -> %q\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	return out
}

// Diff computes the difference between an existing env map and a new one.
// existing is the current state; incoming is the desired state.
func Diff(existing, incoming map[string]string) *DiffResult {
	result := &DiffResult{}

	for k, newVal := range incoming {
		oldVal, exists := existing[k]
		if !exists {
			result.Changes = append(result.Changes, Change{
				Key:      k,
				Type:     Added,
				NewValue: newVal,
			})
		} else if oldVal != newVal {
			result.Changes = append(result.Changes, Change{
				Key:      k,
				Type:     Changed,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	for k, oldVal := range existing {
		if _, exists := incoming[k]; !exists {
			result.Changes = append(result.Changes, Change{
				Key:      k,
				Type:     Removed,
				OldValue: oldVal,
			})
		}
	}

	return result
}
