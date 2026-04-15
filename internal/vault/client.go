package vault

import (
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	vc *vaultapi.Client
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	vc, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	vc.SetToken(token)

	return &Client{vc: vc}, nil
}

// GetSecrets reads a KV v2 secret at the given path and returns key/value pairs.
func (c *Client) GetSecrets(mountPath, secretPath string) (map[string]string, error) {
	fullPath := strings.Trim(mountPath, "/") + "/data/" + strings.Trim(secretPath, "/")

	secret, err := c.vc.Logical().Read(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", fullPath, err)
	}

	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", fullPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected secret data format at path %q", fullPath)
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}

	return result, nil
}
