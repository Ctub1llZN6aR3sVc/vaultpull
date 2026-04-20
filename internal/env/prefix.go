package env

import (
	"fmt"
	"strings"
)

// PrefixOptions controls how prefix addition/removal is applied.
type PrefixOptions struct {
	// Add prepends this string to every key.
	Add string
	// Strip removes this string from the start of every key.
	Strip string
	// FailOnConflict returns an error if adding a prefix would produce a duplicate key.
	FailOnConflict bool
	// DryRun returns the result without mutating the input.
	DryRun bool
}

// PrefixResult holds the outcome of a Prefix operation.
type PrefixResult struct {
	Out      map[string]string
	Renamed  []string
	Skipped  []string
	Conflicts []string
}

// Summary returns a human-readable summary of the prefix operation.
func (r PrefixResult) Summary() string {
	parts := []string{}
	if len(r.Renamed) > 0 {
		parts = append(parts, fmt.Sprintf("%d renamed", len(r.Renamed)))
	}
	if len(r.Skipped) > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", len(r.Skipped)))
	}
	if len(r.Conflicts) > 0 {
		parts = append(parts, fmt.Sprintf("%d conflicts", len(r.Conflicts)))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}

// Prefix applies prefix addition or stripping to the keys of secrets.
// If both Add and Strip are set, Strip is applied first.
func Prefix(secrets map[string]string, opts PrefixOptions) (PrefixResult, error) {
	out := make(map[string]string, len(secrets))
	result := PrefixResult{Out: out}

	for k, v := range secrets {
		newKey := k

		if opts.Strip != "" && strings.HasPrefix(k, opts.Strip) {
			newKey = strings.TrimPrefix(newKey, opts.Strip)
		}

		if opts.Add != "" {
			newKey = opts.Add + newKey
		}

		if newKey == k {
			out[k] = v
			result.Skipped = append(result.Skipped, k)
			continue
		}

		if _, exists := out[newKey]; exists {
			result.Conflicts = append(result.Conflicts, newKey)
			if opts.FailOnConflict {
				return PrefixResult{}, fmt.Errorf("prefix conflict: key %q already exists", newKey)
			}
			continue
		}

		out[newKey] = v
		result.Renamed = append(result.Renamed, fmt.Sprintf("%s -> %s", k, newKey))
	}

	if opts.DryRun {
		result.Out = out
	}

	return result, nil
}
