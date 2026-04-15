package vault

import "fmt"

// MockClient is a test double for Client, useful in higher-level unit tests.
type MockClient struct {
	// Secrets maps "mountPath/secretPath" to a set of key/value pairs.
	Secrets map[string]map[string]string
	// Err, if set, is returned by GetSecrets for every call.
	Err error
}

// NewMockClient returns a MockClient pre-populated with the given secrets map.
func NewMockClient(secrets map[string]map[string]string) *MockClient {
	if secrets == nil {
		secrets = make(map[string]map[string]string)
	}
	return &MockClient{Secrets: secrets}
}

// GetSecrets returns the pre-configured secrets for the given path combination.
func (m *MockClient) GetSecrets(mountPath, secretPath string) (map[string]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	key := mountPath + "/" + secretPath
	secrets, ok := m.Secrets[key]
	if !ok {
		return nil, fmt.Errorf("mock: no secrets configured for path %q", key)
	}

	// Return a copy to avoid mutation between calls.
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	return copy, nil
}
