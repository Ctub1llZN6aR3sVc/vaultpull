package env

// DefaultsOptions controls how defaults are applied.
type DefaultsOptions struct {
	// Overwrite replaces existing keys with defaults.
	Overwrite bool
	// DryRun returns the result without modifying dst.
	DryRun bool
}

// DefaultsResult summarises the apply operation.
type DefaultsResult struct {
	Applied []string
	Skipped []string
}

func (r DefaultsResult) Summary() string {
	if len(r.Applied) == 0 {
		return "defaults: no keys applied"
	}
	return fmt.Sprintf("defaults: applied %d key(s), skipped %d", len(r.Applied), len(r.Skipped))
}

// ApplyDefaults merges default values into dst.
// Keys already present in dst are skipped unless Overwrite is set.
func ApplyDefaults(dst, defaults map[string]string, opts DefaultsOptions) (map[string]string, DefaultsResult) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var result DefaultsResult
	for k, v := range defaults {
		if _, exists := out[k]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		result.Applied = append(result.Applied, k)
		if !opts.DryRun {
			out[k] = v
		}
	}

	sort.Strings(result.Applied)
	sort.Strings(result.Skipped)
	return out, result
}
