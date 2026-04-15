package env

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// templateVarPattern matches ${VAR_NAME} or $VAR_NAME style references.
var templateVarPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// RenderTemplate replaces variable references in a template string with values
// from the provided secrets map, falling back to environment variables.
// Returns an error if any referenced variable is not found.
func RenderTemplate(tmpl string, secrets map[string]string) (string, error) {
	var missing []string

	result := templateVarPattern.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := extractKey(match)
		if val, ok := secrets[key]; ok {
			return val
		}
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		missing = append(missing, key)
		return match
	})

	if len(missing) > 0 {
		return "", fmt.Errorf("template: unresolved variables: %s", strings.Join(missing, ", "))
	}

	return result, nil
}

// RenderMap applies RenderTemplate to every value in the secrets map,
// allowing values to reference other keys in the same map.
func RenderMap(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		rendered, err := RenderTemplate(v, secrets)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = rendered
	}
	return out, nil
}

func extractKey(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
