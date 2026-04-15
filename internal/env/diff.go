package env

import "fmt"

// ChangeType represents the kind of change detected.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Change represents a single key-level difference between two env maps.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// DiffResult holds all changes between two env maps.
type DiffResult struct {
	Changes []Change
}

// IsEmpty returns true when no changes were detected.
func (d *DiffResult) IsEmpty() bool {
	return len(d.Changes) == 0
}

// Summary returns a human-readable summary of the diff.
func (d *DiffResult) Summary() string {
	if d.IsEmpty() {
		return "no changes"
	}
	added, removed, changed := 0, 0, 0
	for _, c := range d.Changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed", added, removed, changed)
}

// Diff compares an existing env map (before) with a new one (after)
// and returns the set of changes.
func Diff(before, after map[string]string) *DiffResult {
	result := &DiffResult{}

	for key, newVal := range after {
		oldVal, exists := before[key]
		if !exists {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Added,
				NewValue: newVal,
			})
		} else if oldVal != newVal {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Changed,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	for key, oldVal := range before {
		if _, exists := after[key]; !exists {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Removed,
				OldValue: oldVal,
			})
		}
	}

	return result
}
