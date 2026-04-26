package env

import (
	"fmt"
	"strings"
)

// SuffixOptions controls how suffixes are added or stripped from keys.
type SuffixOptions struct {
	// Add appends this string to all (or selected) keys.
	Add string
	// Strip removes this string from the end of all (or selected) keys.
	Strip string
	// Keys restricts the operation to a specific set of keys.
	// If empty, all keys are affected.
	Keys []string
	// FailOnConflict returns an error if adding a suffix would produce a key
	// that already exists in the map.
	FailOnConflict bool
	// DryRun returns the result without modifying the original map.
	DryRun bool
}

// SuffixResult holds the output of a Suffix operation.
type SuffixResult struct {
	Secrets   map[string]string
	Affected  []string
	Conflicts []string
}

func (r SuffixResult) Summary() string {
	if len(r.Affected) == 0 {
		return "suffix: no keys affected"
	}
	return fmt.Sprintf("suffix: %d key(s) affected", len(r.Affected))
}

// Suffix adds or strips a suffix from env keys.
func Suffix(secrets map[string]string, opts SuffixOptions) (SuffixResult, error) {
	target := make(map[string]string, len(secrets))
	for k, v := range secrets {
		target[k] = v
	}

	selected := toSuffixSet(opts.Keys)

	result := map[string]string{}
	var affected []string
	var conflicts []string

	for k, v := range target {
		newKey := k

		if len(selected) == 0 || selected[k] {
			if opts.Strip != "" && strings.HasSuffix(k, opts.Strip) {
				newKey = k[:len(k)-len(opts.Strip)]
			}
			if opts.Add != "" {
				newKey = newKey + opts.Add
			}
		}

		if newKey != k {
			if _, exists := target[newKey]; exists {
				conflicts = append(conflicts, newKey)
				if opts.FailOnConflict {
					return SuffixResult{}, fmt.Errorf("suffix: key conflict: %q already exists", newKey)
				}
				result[k] = v
				continue
			}
			affected = append(affected, k)
			result[newKey] = v
		} else {
			result[k] = v
		}
	}

	if !opts.DryRun {
		for k := range secrets {
			delete(secrets, k)
		}
		for k, v := range result {
			secrets[k] = v
		}
	}

	return SuffixResult{Secrets: result, Affected: affected, Conflicts: conflicts}, nil
}

func toSuffixSet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
