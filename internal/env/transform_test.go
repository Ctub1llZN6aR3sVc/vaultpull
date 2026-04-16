package env

import (
	"testing"
)

func TestTransform_NoOptions(t *testing.T) {
	input := map[string]string{"db_host": "localhost", "db_port": "5432"}
	out := Transform(input, TransformOptions{})
	if out["db_host"] != "localhost" || out["db_port"] != "5432" {
		t.Fatal("expected keys to be unchanged")
	}
}

func TestTransform_Uppercase(t *testing.T) {
	input := map[string]string{"db_host": "localhost"}
	out := Transform(input, TransformOptions{Uppercase: true})
	if _, ok := out["DB_HOST"]; !ok {
		t.Fatal("expected key to be uppercased")
	}
}

func TestTransform_AddPrefix(t *testing.T) {
	input := map[string]string{"HOST": "localhost"}
	out := Transform(input, TransformOptions{Prefix: "APP_"})
	if _, ok := out["APP_HOST"]; !ok {
		t.Fatal("expected prefix to be added")
	}
}

func TestTransform_StripPrefix(t *testing.T) {
	input := map[string]string{"SECRET_KEY": "abc"}
	out := Transform(input, TransformOptions{StripPrefix: "SECRET_"})
	if _, ok := out["KEY"]; !ok {
		t.Fatal("expected prefix to be stripped")
	}
}

func TestTransform_StripAndUppercaseAndPrefix(t *testing.T) {
	input := map[string]string{"raw_token": "xyz"}
	out := Transform(input, TransformOptions{
		StripPrefix: "raw_",
		Uppercase:   true,
		Prefix:      "APP_",
	})
	if out["APP_TOKEN"] != "xyz" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestTransform_EmptyKeyAfterStrip_Skipped(t *testing.T) {
	input := map[string]string{"prefix": "value"}
	out := Transform(input, TransformOptions{StripPrefix: "prefix"})
	if len(out) != 0 {
		t.Fatal("expected empty key to be skipped")
	}
}

func TestTransform_PreservesValues(t *testing.T) {
	input := map[string]string{"KEY": "my secret value"}
	out := Transform(input, TransformOptions{Uppercase: true})
	if out["KEY"] != "my secret value" {
		t.Fatal("expected value to be preserved")
	}
}
