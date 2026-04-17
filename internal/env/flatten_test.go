package env

import (
	"testing"
)

func TestFlatten_FlatMap(t *testing.T) {
	input := map[string]any{
		"FOO": "bar",
		"BAZ": "qux",
	}
	out, err := Flatten(input, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestFlatten_NestedMap(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
	}
	out, err := Flatten(input, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %q", out["db_host"])
	}
	if out["db_port"] != "5432" {
		t.Errorf("expected db_port=5432, got %q", out["db_port"])
	}
}

func TestFlatten_Uppercase(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{
			"host": "localhost",
		},
	}
	out, err := Flatten(input, FlattenOptions{Uppercase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out)
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	input := map[string]any{
		"app": map[string]any{
			"name": "vaultpull",
		},
	}
	out, err := Flatten(input, FlattenOptions{Separator: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["app.name"] != "vaultpull" {
		t.Errorf("expected app.name=vaultpull, got %v", out)
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	input := map[string]any{
		"host": "localhost",
	}
	out, err := Flatten(input, FlattenOptions{Prefix: "DB"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_host"] != "localhost" {
		t.Errorf("expected DB_host=localhost, got %v", out)
	}
}

func TestFlatten_NilValue(t *testing.T) {
	input := map[string]any{
		"KEY": nil,
	}
	out, err := Flatten(input, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["KEY"]; !ok || v != "" {
		t.Errorf("expected KEY=empty string, got %q", v)
	}
}
