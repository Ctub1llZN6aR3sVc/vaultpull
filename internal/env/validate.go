package env

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation issue.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) Add(key, message string) {
	r.Errors = append(r.Errors, ValidationError{Key: key, Message: message})
}

func (r *ValidationResult) Summary() string {
	if r.IsValid() {
		return "all secrets valid"
	}
	lines := make([]string, 0, len(r.Errors))
	for _, e := range r.Errors {
		lines = append(lines, "  - "+e.Error())
	}
	return fmt.Sprintf("%d validation error(s):\n%s", len(r.Errors), strings.Join(lines, "\n"))
}

// Validate checks a map of secrets for common issues such as empty keys,
// empty values for sensitive fields, or keys containing invalid characters.
func Validate(secrets map[string]string) ValidationResult {
	result := ValidationResult{}
	for k, v := range secrets {
		if k == "" {
			result.Add(k, "key must not be empty")
			continue
		}
		if strings.ContainsAny(k, " \t\n") {
			result.Add(k, "key must not contain whitespace")
		}
		if IsSensitive(k) && strings.TrimSpace(v) == "" {
			result.Add(k, "sensitive key must not have an empty value")
		}
	}
	return result
}
