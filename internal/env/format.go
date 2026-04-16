package env

import (
	"fmt"
	"sort"
	"strings"
)

// Format controls the output format for env file serialization.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Serialize converts a secrets map into the given format string.
func Serialize(secrets map[string]string, format Format) (string, error) {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch format {
	case FormatExport:
		var sb strings.Builder
		for _, k := range keys {
			fmt.Fprintf(&sb, "export %s=%s\n", k, quoteIfNeeded(secrets[k]))
		}
		return sb.String(), nil

	case FormatJSON:
		var sb strings.Builder
		sb.WriteString("{\n")
		for i, k := range keys {
			comma := ","
			if i == len(keys)-1 {
				comma = ""
			}
			fmt.Fprintf(&sb, "  %q: %q%s\n", k, secrets[k], comma)
		}
		sb.WriteString("}\n")
		return sb.String(), nil

	case FormatDotenv, "":
		var sb strings.Builder
		for _, k := range keys {
			fmt.Fprintf(&sb, "%s=%s\n", k, quoteIfNeeded(secrets[k]))
		}
		return sb.String(), nil

	default:
		return "", fmt.Errorf("unsupported format: %q", format)
	}
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n#$\"") {
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		return `"` + escaped + `"`
	}
	return v
}
