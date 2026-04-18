package env

import "strings"

// TrimOptions controls how secrets are trimmed.
type TrimOptions struct {
	// TrimKeys removes leading/trailing whitespace from keys.
	TrimKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// TrimPrefix removes a prefix from all keys (after optional whitespace trim).
	TrimPrefix string
	// TrimSuffix removes a suffix from all keys.
	TrimSuffix string
}

// Trim returns a new map with keys and/or values trimmed according to opts.
// Keys that become empty after trimming are dropped.
func Trim(secrets map[string]string, opts TrimOptions) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if opts.TrimKeys {
			k = strings.TrimSpace(k)
		}
		if opts.TrimPrefix != "" {
			k = strings.TrimPrefix(k, opts.TrimPrefix)
		}
		if opts.TrimSuffix != "" {
			k = strings.TrimSuffix(k, opts.TrimSuffix)
		}
		if k == "" {
			continue
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		out[k] = v
	}
	return out
}
