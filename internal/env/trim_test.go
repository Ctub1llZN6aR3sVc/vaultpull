package env

import (
	"testing"
)

func TestTrim_NoOptions(t *testing.T) {
	secrets := map[string]string{" KEY ": " value "}
	result := Trim(secrets, TrimOptions{})
	if _, ok := result[" KEY "]; !ok {
		t.Fatal("expected key to be preserved")
	}
	if result[" KEY "] != " value " {
		t.Fatal("expected value to be preserved")
	}
}

func TestTrim_TrimKeys(t *testing.T) {
	secrets := map[string]string{" MY_KEY ": "val"}
	result := Trim(secrets, TrimOptions{TrimKeys: true})
	if _, ok := result["MY_KEY"]; !ok {
		t.Fatalf("expected trimmed key MY_KEY, got %v", result)
	}
}

func TestTrim_TrimValues(t *testing.T) {
	secrets := map[string]string{"KEY": "  hello  "}
	result := Trim(secrets, TrimOptions{TrimValues: true})
	if result["KEY"] != "hello" {
		t.Fatalf("expected trimmed value, got %q", result["KEY"])
	}
}

func TestTrim_TrimPrefix(t *testing.T) {
	secrets := map[string]string{"APP_FOO": "1", "APP_BAR": "2", "OTHER": "3"}
	result := Trim(secrets, TrimOptions{TrimPrefix: "APP_"})
	if _, ok := result["FOO"]; !ok {
		t.Fatal("expected FOO")
	}
	if _, ok := result["BAR"]; !ok {
		t.Fatal("expected BAR")
	}
	if _, ok := result["OTHER"]; !ok {
		t.Fatal("expected OTHER unchanged")
	}
}

func TestTrim_TrimSuffix(t *testing.T) {
	secrets := map[string]string{"KEY_DEV": "v", "NAME_DEV": "n"}
	result := Trim(secrets, TrimOptions{TrimSuffix: "_DEV"})
	if _, ok := result["KEY"]; !ok {
		t.Fatal("expected KEY")
	}
	if _, ok := result["NAME"]; !ok {
		t.Fatal("expected NAME")
	}
}

func TestTrim_EmptyKeyDropped(t *testing.T) {
	secrets := map[string]string{"APP_": "val"}
	result := Trim(secrets, TrimOptions{TrimPrefix: "APP_"})
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %v", result)
	}
}

func TestTrim_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{" K ": " V "}
	_ = Trim(secrets, TrimOptions{TrimKeys: true, TrimValues: true})
	if _, ok := secrets[" K "]; !ok {
		t.Fatal("input map was mutated")
	}
}
