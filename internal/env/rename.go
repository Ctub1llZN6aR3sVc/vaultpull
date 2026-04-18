package env

import "fmt"

// RenameOptions controls how keys are renamed.
type RenameOptions struct {
	// Map is a set of old->new key renames to apply.
	Map map[string]string
	// FailOnMissing returns an error if a key in Map is not found in secrets.
	FailOnMissing bool
}

// RenameResult holds the outcome of a Rename operation.
type RenameResult struct {
	Renamed  []string
	Missed   []string
}

func (r RenameResult) Summary() string {
	if len(r.Renamed) == 0 && len(r.Missed) == 0 {
		return "no renames applied"
	}
	return fmt.Sprintf("%d renamed, %d not found", len(r.Renamed), len(r.Missed))
}

// Rename applies key renames to secrets, returning a new map and a result summary.
// Original keys are removed; values are carried over to new keys.
func Rename(secrets map[string]string, opts RenameOptions) (map[string]string, RenameResult, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var result RenameResult
	for oldKey, newKey := range opts.Map {
		val, ok := out[oldKey]
		if !ok {
			result.Missed = append(result.Missed, oldKey)
			if opts.FailOnMissing {
				return nil, result, fmt.Errorf("rename: key %q not found in secrets", oldKey)
			}
			continue
		}
		delete(out, oldKey)
		out[newKey] = val
		result.Renamed = append(result.Renamed, oldKey)
	}
	return out, result, nil
}
