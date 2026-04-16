package env

import "strings"

// RedactOptions controls which keys are redacted and how.
type RedactOptions struct {
	// Keys to always redact regardless of name.
	Keys []string
	// If true, apply IsSensitive heuristic in addition to explicit Keys.
	AutoDetect bool
	// Placeholder replaces the value; defaults to "[REDACTED]".
	Placeholder string
}

// Redact returns a copy of secrets with sensitive values replaced by a
// placeholder string. It never mutates the input map.
func Redact(secrets map[string]string, opts RedactOptions) map[string]string {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}

	explicit := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[strings.ToUpper(k)] = struct{}{}
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		should := false
		if _, ok := explicit[strings.ToUpper(k)]; ok {
			should = true
		}
		if !should && opts.AutoDetect && IsSensitive(k) {
			should = true
		}
		if should {
			out[k] = placeholder
		} else {
			out[k] = v
		}
	}
	return out
}
