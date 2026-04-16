package env

import (
	"testing"
)

var testKey = []byte("12345678901234567890123456789012") // 32 bytes

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	plaintext := "super-secret-value"
	enc, err := Encrypt([]byte(plaintext), testKey)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if enc == plaintext {
		t.Fatal("encrypted value should differ from plaintext")
	}
	dec, err := Decrypt(enc, testKey)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(dec) != plaintext {
		t.Fatalf("expected %q, got %q", plaintext, dec)
	}
}

func TestEncrypt_DifferentCiphertexts(t *testing.T) {
	plaintext := []byte("value")
	a, _ := Encrypt(plaintext, testKey)
	b, _ := Encrypt(plaintext, testKey)
	if a == b {
		t.Fatal("two encryptions of same plaintext should differ (random nonce)")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("not-base64!!!", testKey)
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecrypt_TruncatedData(t *testing.T) {
	_, err := Decrypt("dG9vc2hvcnQ=", testKey) // "tooshort"
	if err == nil {
		t.Fatal("expected error for truncated ciphertext")
	}
}

func TestEncryptSecrets_Roundtrip(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
	}
	enc, err := EncryptSecrets(secrets, testKey)
	if err != nil {
		t.Fatalf("EncryptSecrets: %v", err)
	}
	dec, err := DecryptSecrets(enc, testKey)
	if err != nil {
		t.Fatalf("DecryptSecrets: %v", err)
	}
	for k, want := range secrets {
		if got := dec[k]; got != want {
			t.Errorf("key %s: expected %q, got %q", k, want, got)
		}
	}
}

func TestDecryptSecrets_BadValue(t *testing.T) {
	secrets := map[string]string{"KEY": "not-valid-ciphertext"}
	_, err := DecryptSecrets(secrets, testKey)
	if err == nil {
		t.Fatal("expected error for invalid ciphertext")
	}
}
