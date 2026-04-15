package sync

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

// SecretClient is the interface for fetching secrets from Vault.
type SecretClient interface {
	GetSecrets(path string) (map[string]string, error)
}

// Syncer orchestrates fetching secrets and writing them to an env file.
type Syncer struct {
	client  SecretClient
	envFile string
}

// New creates a Syncer with the given client and env file path.
func New(client SecretClient, envFile string) *Syncer {
	return &Syncer{client: client, envFile: envFile}
}

// NewFromConfig constructs a Syncer from a loaded config and profile name.
func NewFromConfig(cfg *config.Config, profile string, token string) (*Syncer, string, error) {
	p, err := cfg.GetProfile(profile)
	if err != nil {
		return nil, "", fmt.Errorf("profile %q not found: %w", profile, err)
	}

	client, err := vault.NewClient(cfg.Vault.Address, token)
	if err != nil {
		return nil, "", fmt.Errorf("vault client error: %w", err)
	}

	return New(client, p.EnvFile), p.EnvFile, nil
}

// Run fetches secrets from all given paths, merges them, and writes to the env file.
// It returns a DiffResult describing what changed relative to the existing env file.
func (s *Syncer) Run(paths []string) (*env.DiffResult, error) {
	incoming := make(map[string]string)

	for _, path := range paths {
		secrets, err := s.client.GetSecrets(path)
		if err != nil {
			return nil, fmt.Errorf("failed to get secrets from %q: %w", path, err)
		}
		for k, v := range secrets {
			incoming[k] = v
		}
	}

	existing := env.Read(s.envFile)
	diff := env.Diff(existing, incoming)

	w := env.NewWriter(s.envFile)
	if err := w.Merge(incoming); err != nil {
		return nil, fmt.Errorf("failed to write env file: %w", err)
	}

	return diff, nil
}
