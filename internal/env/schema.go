package env

import (
	"fmt"
	"strings"
)

// SchemaField describes an expected secret key with optional constraints.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  string // expected value prefix/suffix hint for docs
}

// SchemaResult holds the outcome of a schema check.
type SchemaResult struct {
	Missing  []string
	Extra    []string
	HasError bool
}

func (r SchemaResult) Summary() string {
	var sb strings.Builder
	if len(r.Missing) == 0 && len(r.Extra) == 0 {
		sb.WriteString("schema OK: all required keys present")
		return sb.String()
	}
	for _, k := range r.Missing {
		sb.WriteString(fmt.Sprintf("missing required key: %s\n", k))
	}
	for _, k := range r.Extra {
		sb.WriteString(fmt.Sprintf("unexpected key: %s\n", k))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// ValidateSchema checks secrets against a list of schema fields.
// It reports missing required keys and unexpected keys when strict is true.
func ValidateSchema(secrets map[string]string, fields []SchemaField, strict bool) SchemaResult {
	result := SchemaResult{}

	expected := make(map[string]bool, len(fields))
	for _, f := range fields {
		expected[f.Key] = true
		if f.Required {
			if _, ok := secrets[f.Key]; !ok {
				result.Missing = append(result.Missing, f.Key)
				result.HasError = true
			}
		}
	}

	if strict {
		for k := range secrets {
			if !expected[k] {
				result.Extra = append(result.Extra, k)
			}
		}
	}

	return result
}

// EmptyValueKeys returns the keys from secrets that are present but have an empty value.
// This is useful for detecting keys that exist in Vault but were not populated.
func EmptyValueKeys(secrets map[string]string) []string {
	var empty []string
	for k, v := range secrets {
		if strings.TrimSpace(v) == "" {
			empty = append(empty, k)
		}
	}
	return empty
}
