package env

import (
	"fmt"
	"strconv"
	"strings"
)

// ComputeOptions controls how computed keys are derived.
type ComputeOptions struct {
	// Expressions maps new key names to expression strings.
	// Supported ops: "KEY1 + KEY2" (concat), "KEY1 - KEY2" (numeric sub),
	// "KEY1 * KEY2" (numeric mul), "len(KEY)" (string length as string).
	Expressions map[string]string
	// Overwrite allows replacing an existing key.
	Overwrite bool
	// DryRun returns the result without mutating dst.
	DryRun bool
}

// ComputeResult holds the outcome of a Compute call.
type ComputeResult struct {
	Added   []string
	Skipped []string
	Errors  []string
}

func (r ComputeResult) Summary() string {
	if len(r.Errors) > 0 {
		return fmt.Sprintf("compute: %d added, %d skipped, %d errors", len(r.Added), len(r.Skipped), len(r.Errors))
	}
	return fmt.Sprintf("compute: %d added, %d skipped", len(r.Added), len(r.Skipped))
}

// Compute derives new keys from expressions over existing secrets.
func Compute(secrets map[string]string, opts *ComputeOptions) (map[string]string, ComputeResult, error) {
	result := ComputeResult{}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	if opts == nil || len(opts.Expressions) == 0 {
		return out, result, nil
	}

	for newKey, expr := range opts.Expressions {
		if _, exists := out[newKey]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, newKey)
			continue
		}
		val, err := evalExpr(expr, out)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", newKey, err))
			continue
		}
		if !opts.DryRun {
			out[newKey] = val
		}
		result.Added = append(result.Added, newKey)
	}

	if len(result.Errors) > 0 {
		return out, result, fmt.Errorf("compute: %s", strings.Join(result.Errors, "; "))
	}
	return out, result, nil
}

func evalExpr(expr string, secrets map[string]string) (string, error) {
	expr = strings.TrimSpace(expr)

	// len(KEY)
	if strings.HasPrefix(expr, "len(") && strings.HasSuffix(expr, ")") {
		key := expr[4 : len(expr)-1]
		v, ok := secrets[key]
		if !ok {
			return "", fmt.Errorf("key %q not found", key)
		}
		return strconv.Itoa(len(v)), nil
	}

	for _, op := range []string{" + ", " - ", " * "} {
		idx := strings.Index(expr, op)
		if idx < 0 {
			continue
		}
		leftKey := strings.TrimSpace(expr[:idx])
		rightKey := strings.TrimSpace(expr[idx+len(op):])
		leftVal, ok := secrets[leftKey]
		if !ok {
			return "", fmt.Errorf("key %q not found", leftKey)
		}
		rightVal, ok := secrets[rightKey]
		if !ok {
			return "", fmt.Errorf("key %q not found", rightKey)
		}
		switch op {
		case " + ":
			return leftVal + rightVal, nil
		case " - ", " * ":
			l, err := strconv.ParseFloat(leftVal, 64)
			if err != nil {
				return "", fmt.Errorf("key %q is not numeric: %v", leftKey, err)
			}
			r, err := strconv.ParseFloat(rightVal, 64)
			if err != nil {
				return "", fmt.Errorf("key %q is not numeric: %v", rightKey, err)
			}
			var res float64
			if op == " - " {
				res = l - r
			} else {
				res = l * r
			}
			return strconv.FormatFloat(res, 'f', -1, 64), nil
		}
	}
	return "", fmt.Errorf("unsupported expression: %q", expr)
}
