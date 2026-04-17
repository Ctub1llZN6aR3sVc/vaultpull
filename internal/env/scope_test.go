package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyScopes_PartitionsCorrectly(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
		"LOG_LEVEL":   "info",
	}
	scopes := []Scope{
		{Name: "database", Keys: []string{"DB_HOST", "DB_PASSWORD"}},
		{Name: "api", Keys: []string{"API_KEY"}},
	}

	result := ApplyScopes(secrets, scopes)

	assert.Equal(t, "localhost", result.Scoped["database"]["DB_HOST"])
	assert.Equal(t, "secret", result.Scoped["database"]["DB_PASSWORD"])
	assert.Equal(t, "abc123", result.Scoped["api"]["API_KEY"])
	assert.Equal(t, map[string]string{"LOG_LEVEL": "info"}, result.Unscoped)
}

func TestApplyScopes_AllUnscoped(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := ApplyScopes(secrets, []Scope{})

	assert.Empty(t, result.Scoped)
	assert.Equal(t, secrets, result.Unscoped)
}

func TestApplyScopes_MissingKeyInScope(t *testing.T) {
	secrets := map[string]string{"PRESENT": "yes"}
	scopes := []Scope{
		{Name: "group", Keys: []string{"PRESENT", "ABSENT"}},
	}

	result := ApplyScopes(secrets, scopes)

	assert.Equal(t, map[string]string{"PRESENT": "yes"}, result.Scoped["group"])
	assert.Empty(t, result.Unscoped)
}

func TestApplyScopes_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"X": "1"}
	original := map[string]string{"X": "1"}
	scopes := []Scope{{Name: "s", Keys: []string{"X"}}}

	ApplyScopes(secrets, scopes)

	assert.Equal(t, original, secrets)
}

func TestScopeNames_ReturnsSorted(t *testing.T) {
	r := ScopeResult{
		Scoped: map[string]map[string]string{
			"zebra": {},
			"alpha": {},
			"mango": {},
		},
	}
	names := ScopeNames(r)
	assert.Equal(t, []string{"alpha", "mango", "zebra"}, names)
}
