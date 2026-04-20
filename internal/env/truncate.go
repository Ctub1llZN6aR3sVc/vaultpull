package env

import "fmt"

// TruncateOptions controls how values are truncated.
type TruncateOptions struct {
	// MaxLength is the maximum allowed value length. Zero means no limit.
	MaxLength int
	// Suffix is appended to truncated values (e.g. "..."). Defaults to "...".
	Suffix string
	// Keys restricts truncation to specific keys. Empty means all keys.
	Keys []string
	// DryRun returns what would change without mutating.
	DryRun bool
}

// TruncateResult holds the outcome of a Truncate call.
type TruncateResult struct {
	Truncated []string
}

// Summary returns a human-readable summary.
func (r TruncateResult) Summary() string {
	if len(r.Truncated) == 0 {
		return "truncate: no values exceeded the limit"
	}
	return fmt.Sprintf("truncate: %d value(s) truncated: %v", len(r.Truncated), r.Truncated)
}

// Truncate shortens values that exceed MaxLength, appending Suffix.
// It returns a new map and a result summary; the input is never mutated.
func Truncate(secrets map[string]string, opts TruncateOptions) (map[string]string, TruncateResult) {
	if opts.MaxLength <= 0 {
		return copyMap(secrets), TruncateResult{}
	}

	suffix := opts.Suffix
	if suffix == "" {
		suffix = "..."
	}

	keySet := toSet(opts.Keys)
	out := make(map[string]string, len(secrets))
	var result TruncateResult

	for k, v := range secrets {
		if len(keySet) > 0 && !keySet[k] {
			out[k] = v
			continue
		}
		if len(v) > opts.MaxLength {
			result.Truncated = append(result.Truncated, k)
			if !opts.DryRun {
				cutAt := opts.MaxLength - len(suffix)
				if cutAt < 0 {
					cutAt = 0
				}
				out[k] = v[:cutAt] + suffix
				continue
			}
		}
		out[k] = v
	}

	return out, result
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
