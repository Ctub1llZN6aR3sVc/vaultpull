package env

import (
	"regexp"
	"strings"
)

// FilterOptions controls which keys are included or excluded.
type FilterOptions struct {
	// IncludeKeys is a list of exact key names to include. If empty, all keys are included.
	IncludeKeys []string
	// ExcludeKeys is a list of exact key names to exclude.
	ExcludeKeys []string
	// IncludePrefix filters keys that start with any of these prefixes.
	IncludePrefix []string
	// ExcludePrefix filters out keys that start with any of these prefixes.
	ExcludePrefix []string
	// IncludePattern is an optional regex; only matching keys are included.
	IncludePattern string
}

// Filter returns a filtered copy of secrets based on FilterOptions.
func Filter(secrets map[string]string, opts FilterOptions) (map[string]string, error) {
	var includeRe *regexp.Regexp
	if opts.IncludePattern != "" {
		var err error
		includeRe, err = regexp.Compile(opts.IncludePattern)
		if err != nil {
			return nil, err
		}
	}

	excludeSet := toSet(opts.ExcludeKeys)
	includeSet := toSet(opts.IncludeKeys)

	result := make(map[string]string)
	for k, v := range secrets {
		if len(includeSet) > 0 && !includeSet[k] {
			continue
		}
		if excludeSet[k] {
			continue
		}
		if len(opts.IncludePrefix) > 0 && !hasAnyPrefix(k, opts.IncludePrefix) {
			continue
		}
		if len(opts.ExcludePrefix) > 0 && hasAnyPrefix(k, opts.ExcludePrefix) {
			continue
		}
		if includeRe != nil && !includeRe.MatchString(k) {
			continue
		}
		result[k] = v
	}
	return result, nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
