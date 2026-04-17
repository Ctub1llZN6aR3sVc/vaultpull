package env

import "sort"

// Scope represents a named grouping of secret keys.
type Scope struct {
	Name string
	Keys []string
}

// ScopeResult holds the result of scoping secrets.
type ScopeResult struct {
	Scoped   map[string]map[string]string // scope name -> secrets
	Unscoped map[string]string            // keys not matched by any scope
}

// ApplyScopes partitions secrets into named scopes based on key membership.
// Keys not matched by any scope are placed in Unscoped.
func ApplyScopes(secrets map[string]string, scopes []Scope) ScopeResult {
	result := ScopeResult{
		Scoped:   make(map[string]map[string]string),
		Unscoped: make(map[string]string),
	}

	assigned := make(map[string]bool)

	for _, scope := range scopes {
		bucket := make(map[string]string)
		for _, key := range scope.Keys {
			if val, ok := secrets[key]; ok {
				bucket[key] = val
				assigned[key] = true
			}
		}
		result.Scoped[scope.Name] = bucket
	}

	for k, v := range secrets {
		if !assigned[k] {
			result.Unscoped[k] = v
		}
	}

	return result
}

// ScopeNames returns sorted scope names from a ScopeResult.
func ScopeNames(r ScopeResult) []string {
	names := make([]string, 0, len(r.Scoped))
	for name := range r.Scoped {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
