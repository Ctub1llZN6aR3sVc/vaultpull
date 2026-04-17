package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare_IdenticalMaps(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Compare(a, b)
	assert.True(t, r.IsClean())
	assert.ElementsMatch(t, []string{"FOO", "BAZ"}, r.Identical)
}

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "x"}
	b := map[string]string{"FOO": "bar"}
	r := Compare(a, b)
	assert.Equal(t, []string{"ONLY_A"}, r.OnlyInA)
	assert.Empty(t, r.OnlyInB)
	assert.False(t, r.IsClean())
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "y"}
	r := Compare(a, b)
	assert.Equal(t, []string{"ONLY_B"}, r.OnlyInB)
	assert.Empty(t, r.OnlyInA)
}

func TestCompare_Differ(t *testing.T) {
	a := map[string]string{"KEY": "old"}
	b := map[string]string{"KEY": "new"}
	r := Compare(a, b)
	assert.Equal(t, []string{"KEY"}, r.Differ)
	assert.Empty(t, r.Identical)
	assert.False(t, r.IsClean())
}

func TestCompare_Mixed(t *testing.T) {
	a := map[string]string{"SAME": "v", "CHANGED": "old", "GONE": "x"}
	b := map[string]string{"SAME": "v", "CHANGED": "new", "NEW": "y"}
	r := Compare(a, b)
	assert.Equal(t, []string{"GONE"}, r.OnlyInA)
	assert.Equal(t, []string{"NEW"}, r.OnlyInB)
	assert.Equal(t, []string{"CHANGED"}, r.Differ)
	assert.Equal(t, []string{"SAME"}, r.Identical)
}

func TestCompare_EmptyMaps(t *testing.T) {
	r := Compare(map[string]string{}, map[string]string{})
	assert.True(t, r.IsClean())
	assert.Empty(t, r.Identical)
}
