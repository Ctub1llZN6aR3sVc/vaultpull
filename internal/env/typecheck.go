package env

import (
	"fmt"
	"strconv"
	"strings"
)

// TypeRule defines an expected type for a secret key.
type TypeRule struct {
	Key      string
	Expected string // "int", "float", "bool", "string"
}

// TypeCheckResult holds the outcome of a type check operation.
type TypeCheckResult struct {
	Violations []TypeViolation
}

// TypeViolation describes a single type mismatch.
type TypeViolation struct {
	Key      string
	Expected string
	Actual   string
	Value    string
}

func (r TypeCheckResult) OK() bool {
	return len(r.Violations) == 0
}

func (r TypeCheckResult) Summary() string {
	if r.OK() {
		return "typecheck: all values match expected types"
	}
	var sb strings.Builder
	for _, v := range r.Violations {
		fmt.Fprintf(&sb, "typecheck: key %q expected %s, got %q\n", v.Key, v.Expected, v.Value)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// TypeCheck validates that secret values conform to the provided type rules.
func TypeCheck(secrets map[string]string, rules []TypeRule) TypeCheckResult {
	var violations []TypeViolation
	for _, rule := range rules {
		val, ok := secrets[rule.Key]
		if !ok {
			continue
		}
		if err := checkType(val, rule.Expected); err != nil {
			violations = append(violations, TypeViolation{
				Key:      rule.Key,
				Expected: rule.Expected,
				Actual:   inferType(val),
				Value:    val,
			})
		}
	}
	return TypeCheckResult{Violations: violations}
}

func checkType(val, expected string) error {
	switch strings.ToLower(expected) {
	case "int":
		_, err := strconv.ParseInt(val, 10, 64)
		return err
	case "float":
		_, err := strconv.ParseFloat(val, 64)
		return err
	case "bool":
		_, err := strconv.ParseBool(val)
		return err
	case "string":
		return nil
	default:
		return fmt.Errorf("unknown type %q", expected)
	}
}

func inferType(val string) string {
	if _, err := strconv.ParseInt(val, 10, 64); err == nil {
		return "int"
	}
	if _, err := strconv.ParseFloat(val, 64); err == nil {
		return "float"
	}
	if _, err := strconv.ParseBool(val); err == nil {
		return "bool"
	}
	return "string"
}
