package env

import (
	"strings"
	"testing"
)

func TestMaskWrite_NoOptions(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "PASSWORD": "secret"}
	out, res := MaskWrite(secrets, MaskWriteOptions{})
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", out["FOO"])
	}
	if out["PASSWORD"] != "secret" {
		t.Errorf("expected PASSWORD=secret, got %s", out["PASSWORD"])
	}
	if len(res.Masked) != 0 {
		t.Errorf("expected no masked keys, got %v", res.Masked)
	}
}

func TestMaskWrite_AutoDetect(t *testing.T) {
	secrets := map[string]string{"API_KEY": "abc123", "HOST": "localhost"}
	out, res := MaskWrite(secrets, MaskWriteOptions{AutoDetect: true})
	if out["API_KEY"] != "***" {
		t.Errorf("expected API_KEY masked, got %s", out["API_KEY"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", out["HOST"])
	}
	if len(res.Masked) != 1 || res.Masked[0] != "API_KEY" {
		t.Errorf("unexpected masked keys: %v", res.Masked)
	}
}

func TestMaskWrite_ExplicitKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, res := MaskWrite(secrets, MaskWriteOptions{Keys: []string{"FOO"}})
	if out["FOO"] != "***" {
		t.Errorf("expected FOO masked")
	}
	if out["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux")
	}
	if len(res.Masked) != 1 {
		t.Errorf("expected 1 masked key")
	}
}

func TestMaskWrite_CustomPlaceholder(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abc"}
	out, _ := MaskWrite(secrets, MaskWriteOptions{Keys: []string{"TOKEN"}, Placeholder: "REDACTED"})
	if out["TOKEN"] != "REDACTED" {
		t.Errorf("expected REDACTED, got %s", out["TOKEN"])
	}
}

func TestMaskWrite_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"PASSWORD": "secret"}
	MaskWrite(secrets, MaskWriteOptions{AutoDetect: true})
	if secrets["PASSWORD"] != "secret" {
		t.Errorf("input was mutated")
	}
}

func TestMaskWrite_Summary(t *testing.T) {
	secrets := map[string]string{"API_KEY": "x", "TOKEN": "y"}
	_, res := MaskWrite(secrets, MaskWriteOptions{Keys: []string{"API_KEY", "TOKEN"}})
	s := res.Summary()
	if !strings.Contains(s, "2 key(s)") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestMaskWrite_SummaryEmpty(t *testing.T) {
	_, res := MaskWrite(map[string]string{}, MaskWriteOptions{})
	if res.Summary() != "mask_write: no keys masked" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}
