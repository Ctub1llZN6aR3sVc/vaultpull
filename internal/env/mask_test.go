package env

import (
	"testing"
)

func TestIsSensitive_DetectsSensitiveKeys(t *testing.T) {
	sensitive := []string{
		"DB_PASSWORD",
		"AWS_SECRET_ACCESS_KEY",
		"GITHUB_TOKEN",
		"PRIVATE_KEY",
		"API_KEY",
		"AUTH_TOKEN",
		"DSN",
		"TLS_CERT",
	}
	for _, key := range sensitive {
		if !IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_AllowsNonSensitiveKeys(t *testing.T) {
	public := []string{
		"APP_ENV",
		"LOG_LEVEL",
		"PORT",
		"HOST",
		"REGION",
	}
	for _, key := range public {
		if IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestMaskValue_ShortValue(t *testing.T) {
	result := MaskValue("abc")
	if result != "****" {
		t.Errorf("expected '****', got %q", result)
	}
}

func TestMaskValue_LongValue(t *testing.T) {
	result := MaskValue("supersecretvalue")
	if result != "supe****" {
		t.Errorf("expected 'supe****', got %q", result)
	}
}

func TestMaskValue_ExactlyFourChars(t *testing.T) {
	result := MaskValue("1234")
	if result != "****" {
		t.Errorf("expected '****' for 4-char value, got %q", result)
	}
}

func TestMaskSecrets_MasksSensitiveOnly(t *testing.T) {
	input := map[string]string{
		"DB_PASSWORD":  "hunter2",
		"APP_ENV":      "production",
		"API_KEY":      "abcdefgh",
		"LOG_LEVEL":    "info",
	}

	result := MaskSecrets(input)

	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be masked, got %q", result["APP_ENV"])
	}
	if result["LOG_LEVEL"] != "info" {
		t.Errorf("LOG_LEVEL should not be masked, got %q", result["LOG_LEVEL"])
	}
	if result["DB_PASSWORD"] == "hunter2" {
		t.Error("DB_PASSWORD should be masked")
	}
	if result["API_KEY"] == "abcdefgh" {
		t.Error("API_KEY should be masked")
	}
}

func TestMaskSecrets_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{
		"DB_PASSWORD": "hunter2",
	}

	_ = MaskSecrets(input)

	if input["DB_PASSWORD"] != "hunter2" {
		t.Error("original map should not be mutated")
	}
}
