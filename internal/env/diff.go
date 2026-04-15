package env

// DiffResult holds the result of comparing two env maps.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]OldNew
}

// OldNew holds the old and new values for a changed key.
type OldNew struct {
	Old string
	New string
}

// Diff compares a current env map with an incoming map and returns
// the keys that were added, removed, or changed.
func Diff(current, incoming map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]OldNew),
	}

	for k, newVal := range incoming {
		oldVal, exists := current[k]
		if !exists {
			result.Added[k] = newVal
		} else if oldVal != newVal {
			result.Changed[k] = OldNew{Old: oldVal, New: newVal}
		}
	}

	for k, oldVal := range current {
		if _, exists := incoming[k]; !exists {
			result.Removed[k] = oldVal
		}
	}

	return result
}

// IsEmpty returns true when there are no differences.
func (d DiffResult) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}
