package env

import (
	"strings"
)

// NormalizeOptions controls how secret keys and values are normalized.
type NormalizeOptions struct {
	UppercaseKeys  bool
	TrimValues     bool
	ReplaceHyphens bool // replace hyphens in keys with underscores
}

// NormalizeResult holds the outcome of a Normalize call.
type NormalizeResult struct {
	Output  map[string]string
	Renamed []string // keys that were altered
}

// Summary returns a human-readable summary of the normalization.
func (r NormalizeResult) Summary() string {
	if len(r.Renamed) == 0 {
		return "normalize: no keys altered"
	}
	return "normalize: " + strings.Join(r.Renamed, ", ") + " altered"
}

// Normalize applies key/value normalization to secrets according to opts.
func Normalize(secrets map[string]string, opts NormalizeOptions) NormalizeResult {
	out := make(map[string]string, len(secrets))
	var renamed []string

	for k, v := range secrets {
		newKey := k
		if opts.ReplaceHyphens {
			newKey = strings.ReplaceAll(newKey, "-", "_")
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if newKey != k {
			renamed = append(renamed, k+"→"+newKey)
		}
		out[newKey] = v
	}
	return NormalizeResult{Output: out, Renamed: renamed}
}
