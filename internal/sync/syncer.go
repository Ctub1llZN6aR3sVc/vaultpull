package sync

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

// VaultClient defines the interface for fetching secrets from Vault.
type VaultClient interface {
	GetSecrets(path string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets from Vault and writing them to .env files.
type Syncer struct {
	client VaultClient
	writer *env.Writer
}

// New creates a new Syncer using the provided VaultClient and env Writer.
func New(client VaultClient, writer *env.Writer) *Syncer {
	return &Syncer{
		client: client,
		writer: writer,
	}
}

// NewFromConfig constructs a Syncer from a resolved config profile.
func NewFromConfig(profile *config.Profile) (*Syncer, error) {
	client, err := vault.NewClient(profile.VaultAddr, profile.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	writer, err := env.NewWriter(profile.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create env writer: %w", err)
	}

	return New(client, writer), nil
}

// Run pulls secrets from each configured Vault path and merges them into the env file.
func (s *Syncer) Run(paths []string) error {
	merged := make(map[string]string)

	for _, path := range paths {
		secrets, err := s.client.GetSecrets(path)
		if err != nil {
			return fmt.Errorf("failed to get secrets from path %q: %w", path, err)
		}
		for k, v := range secrets {
			merged[k] = v
		}
	}

	if err := s.writer.Merge(merged); err != nil {
		return fmt.Errorf("failed to write secrets to env file: %w", err)
	}

	return nil
}
