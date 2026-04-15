package env

import "fmt"

// ChangeType represents the kind of change detected.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
)

// Change describes a single key-level difference between two env maps.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// String returns a human-readable description of the change.
func (c Change) String() string {
	switch c.Type {
	case ChangeAdded:
		return fmt.Sprintf("+ %s", c.Key)
	case ChangeRemoved:
		return fmt.Sprintf("- %s", c.Key)
	case ChangeUpdated:
		return fmt.Sprintf("~ %s", c.Key)
	default:
		return c.Key
	}
}

// DiffResult holds the full set of changes between two env maps.
type DiffResult struct {
	Changes []Change
}

// IsEmpty returns true when no changes were detected.
func (d DiffResult) IsEmpty() bool {
	return len(d.Changes) == 0
}

// Diff compares an existing env map (old) against incoming secrets (new)
// and returns a DiffResult describing what changed.
func Diff(old, incoming map[string]string) DiffResult {
	var changes []Change

	for k, newVal := range incoming {
		oldVal, exists := old[k]
		if !exists {
			changes = append(changes, Change{Key: k, Type: ChangeAdded, NewValue: newVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Type: ChangeUpdated, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, oldVal := range old {
		if _, exists := incoming[k]; !exists {
			changes = append(changes, Change{Key: k, Type: ChangeRemoved, OldValue: oldVal})
		}
	}

	return DiffResult{Changes: changes}
}
