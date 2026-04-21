package env

import (
	"fmt"
	"strings"
)

// RequiredOptions controls behaviour of the Required check.
type RequiredOptions struct {
	// Keys is the explicit list of keys that must be present and non-empty.
	Keys []string
}

// RequiredResult holds the outcome of a Required check.
type RequiredResult struct {
	Missing []string // keys absent from secrets
	Empty   []string // keys present but with empty values
}

// OK returns true when no violations were found.
func (r RequiredResult) OK() bool {
	return len(r.Missing) == 0 && len(r.Empty) == 0
}

// Summary returns a human-readable description of the result.
func (r RequiredResult) Summary() string {
	if r.OK() {
		return "all required keys present"
	}
	var parts []string
	if len(r.Missing) > 0 {
		parts = append(parts, fmt.Sprintf("missing: %s", strings.Join(r.Missing, ", ")))
	}
	if len(r.Empty) > 0 {
		parts = append(parts, fmt.Sprintf("empty: %s", strings.Join(r.Empty, ", ")))
	}
	return strings.Join(parts, "; ")
}

// Required checks that every key listed in opts.Keys exists in secrets and
// has a non-empty value.  It returns a RequiredResult and, when violations
// exist, a non-nil error whose message matches Summary.
func Required(secrets map[string]string, opts RequiredOptions) (RequiredResult, error) {
	var res RequiredResult
	for _, k := range opts.Keys {
		v, ok := secrets[k]
		if !ok {
			res.Missing = append(res.Missing, k)
		} else if strings.TrimSpace(v) == "" {
			res.Empty = append(res.Empty, k)
		}
	}
	if !res.OK() {
		return res, fmt.Errorf("required check failed: %s", res.Summary())
	}
	return res, nil
}
