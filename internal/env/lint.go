package env

import (
	"fmt"
	"strings"
)

// LintIssue represents a single linting problem found in a secrets map.
type LintIssue struct {
	Key     string
	Message string
}

func (l LintIssue) String() string {
	return fmt.Sprintf("%s: %s", l.Key, l.Message)
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) IsClean() bool {
	return len(r.Issues) == 0
}

func (r *LintResult) Summary() string {
	if r.IsClean() {
		return "no lint issues found"
	}
	return fmt.Sprintf("%d lint issue(s) found", len(r.Issues))
}

// Lint checks secrets for common issues such as whitespace in keys,
// duplicate keys (case-insensitive), and overly long values.
func Lint(secrets map[string]string) LintResult {
	result := LintResult{}
	seen := map[string]string{}

	for k, v := range secrets {
		if strings.TrimSpace(k) != k {
			result.Issues = append(result.Issues, LintIssue{Key: k, Message: "key has leading or trailing whitespace"})
		}
		if strings.Contains(k, " ") {
			result.Issues = append(result.Issues, LintIssue{Key: k, Message: "key contains spaces"})
		}
		lower := strings.ToLower(k)
		if prev, ok := seen[lower]; ok && prev != k {
			result.Issues = append(result.Issues, LintIssue{Key: k, Message: fmt.Sprintf("duplicate key (case-insensitive conflict with %q)", prev)})
		} else {
			seen[lower] = k
		}
		if len(v) > 4096 {
			result.Issues = append(result.Issues, LintIssue{Key: k, Message: "value exceeds 4096 characters"})
		}
	}
	return result
}
