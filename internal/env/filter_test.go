package env

import (
	"testing"
)

var baseSecrets = map[string]string{
	"DB_HOST":       "localhost",
	"DB_PASSWORD":   "secret",
	"APP_PORT":      "8080",
	"APP_DEBUG":     "true",
	"INTERNAL_KEY":  "hidden",
	"AWS_ACCESS_KEY": "AKIA123",
}

func TestFilter_NoOptions(t *testing.T) {
	result, err := Filter(baseSecrets, FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(baseSecrets) {
		t.Errorf("expected %d keys, got %d", len(baseSecrets), len(result))
	}
}

func TestFilter_IncludeKeys(t *testing.T) {
	result, err := Filter(baseSecrets, FilterOptions{
		IncludeKeys: []string{"DB_HOST", "APP_PORT"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestFilter_ExcludeKeys(t *testing.T) {
	result, err := Filter(baseSecrets, FilterOptions{
		ExcludeKeys: []string{"DB_PASSWORD", "AWS_ACCESS_KEY"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if _, ok := result["AWS_ACCESS_KEY"]; ok {
		t.Error("AWS_ACCESS_KEY should be excluded")
	}
}

func TestFilter_IncludePrefix(t *testing.T) {
	result, err := Filter(baseSecrets, FilterOptions{
		IncludePrefix: []string{"DB_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys with DB_ prefix, got %d", len(result))
	}
}

func TestFilter_ExcludePrefix(t *testing.T) {
	result, err := Filter(baseSecrets, FilterOptions{
		ExcludePrefix: []string{"INTERNAL_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["INTERNAL_KEY"]; ok {
		t.Error("INTERNAL_KEY should be excluded")
	}
}

func TestFilter_IncludePattern(t *testing.T) {
	result, err := Filter(baseSecrets, FilterOptions{
		IncludePattern: "^APP_",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 APP_ keys, got %d", len(result))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := Filter(baseSecrets, FilterOptions{
		IncludePattern: "[invalid",
	})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestFilter_CombinedPrefixAndExclude(t *testing.T) {
	result, err := Filter(baseSecrets, FilterOptions{
		IncludePrefix: []string{"DB_"},
		ExcludeKeys:   []string{"DB_PASSWORD"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}
