package env

import (
	"strings"
)

// SanitizeOptions controls how secrets are sanitized.
type SanitizeOptions struct {
	// StripControlChars removes non-printable ASCII control characters from values.
	StripControlChars bool
	// TrimQuotes removes surrounding single or double quotes from values.
	TrimQuotes bool
	// NormalizeNewlines replaces \r\n and \r with \n in values.
	NormalizeNewlines bool
}

// SanitizeResult holds the outcome of a sanitize operation.
type SanitizeResult struct {
	Sanitized map[string]string
	ChangedKeys []string
}

// Summary returns a human-readable summary.
func (r SanitizeResult) Summary() string {
	if len(r.ChangedKeys) == 0 {
		return "sanitize: no changes"
	}
	return "sanitize: modified keys: " + strings.Join(r.ChangedKeys, ", ")
}

// Sanitize cleans secret values according to the provided options.
func Sanitize(secrets map[string]string, opts SanitizeOptions) SanitizeResult {
	out := make(map[string]string, len(secrets))
	var changed []string

	for k, v := range secrets {
		original := v

		if opts.NormalizeNewlines {
			v = strings.ReplaceAll(v, "\r\n", "\n")
			v = strings.ReplaceAll(v, "\r", "\n")
		}

		if opts.StripControlChars {
			var b strings.Builder
			for _, r := range v {
				if r == '\n' || r == '\t' || (r >= 0x20 && r != 0x7F) {
					b.WriteRune(r)
				}
			}
			v = b.String()
		}

		if opts.TrimQuotes {
			if len(v) >= 2 {
				if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
					v = v[1 : len(v)-1]
				}
			}
		}

		out[k] = v
		if v != original {
			changed = append(changed, k)
		}
	}

	return SanitizeResult{Sanitized: out, ChangedKeys: changed}
}
