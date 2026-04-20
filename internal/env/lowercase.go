package env

import (
	"strings"
)

// LowercaseOptions controls how Lowercase behaves.
type LowercaseOptions struct {
	// LowercaseKeys converts all key names to lowercase.
	LowercaseKeys bool
	// LowercaseValues converts all values to lowercase.
	LowercaseValues bool
	// OnlyKeys restricts lowercasing to the specified keys (applies to values only).
	OnlyKeys []string
}

// LowercaseResult holds the result of a Lowercase operation.
type LowercaseResult struct {
	Output  map[string]string
	Changed []string
}

// Summary returns a human-readable summary of the lowercase operation.
func (r LowercaseResult) Summary() string {
	if len(r.Changed) == 0 {
		return "lowercase: no changes"
	}
	return "lowercase: changed " + joinLowercaseKeys(r.Changed)
}

func joinLowercaseKeys(keys []string) string {
	if len(keys) == 0 {
		return ""
	}
	out := ""
	for i, k := range keys {
		if i > 0 {
			out += ", "
		}
		out += k
	}
	return out
}

// Lowercase applies lowercase transformations to keys and/or values of secrets.
// It never mutates the input map.
func Lowercase(secrets map[string]string, opts LowercaseOptions) LowercaseResult {
	only := make(map[string]bool, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		only[k] = true
	}

	out := make(map[string]string, len(secrets))
	var changed []string

	for k, v := range secrets {
		newKey := k
		newVal := v

		if opts.LowercaseKeys {
			newKey = strings.ToLower(k)
		}

		if opts.LowercaseValues {
			if len(only) == 0 || only[k] {
				newVal = strings.ToLower(v)
			}
		}

		if newKey != k || newVal != v {
			changed = append(changed, k)
		}

		out[newKey] = newVal
	}

	return LowercaseResult{Output: out, Changed: changed}
}
