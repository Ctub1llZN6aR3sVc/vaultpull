package env

import (
	"fmt"
	"strings"
)

// CoerceOptions controls how values are coerced.
type CoerceOptions struct {
	// BoolValues maps canonical true/false representations.
	NormalizeBools bool
	// TrimWhitespace removes surrounding whitespace from values.
	TrimWhitespace bool
	// LowercaseValues lowercases all values.
	LowercaseValues bool
	// UppercaseValues uppercases all values.
	UppercaseValues bool
}

// CoerceResult holds the outcome of a Coerce operation.
type CoerceResult struct {
	Coerced map[string]string
	Changes []string
}

func (r CoerceResult) Summary() string {
	if len(r.Changes) == 0 {
		return "coerce: no changes"
	}
	return fmt.Sprintf("coerce: %d value(s) changed: %s", len(r.Changes), strings.Join(r.Changes, ", "))
}

var boolTrue = map[string]string{
	"1": "true", "yes": "true", "on": "true", "true": "true",
}
var boolFalse = map[string]string{
	"0": "false", "no": "false", "off": "false", "false": "false",
}

// Coerce applies value normalization rules to a secrets map.
func Coerce(secrets map[string]string, opts CoerceOptions) CoerceResult {
	out := make(map[string]string, len(secrets))
	var changes []string

	for k, v := range secrets {
		orig := v

		if opts.TrimWhitespace {
			v = strings.TrimSpace(v)
		}
		if opts.NormalizeBools {
			lower := strings.ToLower(strings.TrimSpace(v))
			if canonical, ok := boolTrue[lower]; ok {
				v = canonical
			} else if canonical, ok := boolFalse[lower]; ok {
				v = canonical
			}
		}
		if opts.LowercaseValues {
			v = strings.ToLower(v)
		} else if opts.UppercaseValues {
			v = strings.ToUpper(v)
		}

		out[k] = v
		if v != orig {
			changes = append(changes, k)
		}
	}

	return CoerceResult{Coerced: out, Changes: changes}
}
