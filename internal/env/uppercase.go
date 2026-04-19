package env

import (
	"strings"
)

// UppercaseOptions controls the behaviour of Uppercase.
type UppercaseOptions struct {
	// Keys uppercases map keys.
	Keys bool
	// Values uppercases map values.
	Values bool
	// OnlyKeys restricts key uppercasing to this set; empty means all keys.
	OnlyKeys []string
}

// UppercaseResult holds the outcome of an Uppercase call.
type UppercaseResult struct {
	Changed []string
}

func (r UppercaseResult) Summary() string {
	if len(r.Changed) == 0 {
		return "uppercase: no changes"
	}
	return "uppercase: changed " + joinUppercaseKeys(r.Changed)
}

func joinUppercaseKeys(keys []string) string {
	return strings.Join(keys, ", ")
}

// Uppercase returns a new map with keys and/or values uppercased according to
// opts. The original map is never mutated.
func Uppercase(secrets map[string]string, opts UppercaseOptions) (map[string]string, UppercaseResult) {
	allow := toSet(opts.OnlyKeys)
	out := make(map[string]string, len(secrets))
	var changed []string

	for k, v := range secrets {
		newKey := k
		newVal := v

		if opts.Keys {
			if len(allow) == 0 || allow[k] {
				newKey = strings.ToUpper(k)
			}
		}

		if opts.Values {
			if len(allow) == 0 || allow[k] {
				newVal = strings.ToUpper(v)
			}
		}

		if newKey != k || newVal != v {
			changed = append(changed, k)
		}

		out[newKey] = newVal
	}

	return out, UppercaseResult{Changed: changed}
}
