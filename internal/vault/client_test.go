package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/vaultpull/internal/vault"
)

func newMockVaultServer(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestGetSecrets_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
			},
		},
	}

	srv := newMockVaultServer(t, "/v1/secret/data/myapp", payload)
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	secrets, err := client.GetSecrets("secret", "myapp")
	if err != nil {
		t.Fatalf("GetSecrets error: %v", err)
	}

	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
	if secrets["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", secrets["DB_PORT"])
	}
}

func TestGetSecrets_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetSecrets("secret", "missing")
	if err == nil {
		t.Error("expected error for missing secret, got nil")
	}
}

func TestNewClient_InvalidAddress(t *testing.T) {
	_, err := vault.NewClient("://bad-url", "token")
	if err == nil {
		t.Error("expected error for invalid address, got nil")
	}
}
