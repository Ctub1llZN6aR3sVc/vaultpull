package sync

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

// VaultClient is the interface used to fetch secrets.
type VaultClient interface {
	GetSecrets(path string) (map[string]string, error)
}

// Syncer pulls secrets from Vault and writes them to a local .env file.
type Syncer struct {
	client  VaultClient
	envPath string
	paths   []string
	backup  bool
}

// New creates a Syncer with explicit parameters.
func New(client VaultClient, envPath string, paths []string, backup bool) *Syncer {
	return &Syncer{client: client, envPath: envPath, paths: paths, backup: backup}
}

// NewFromConfig builds a Syncer from a resolved config profile.
func NewFromConfig(profile *config.Profile, token string) (*Syncer, error) {
	client, err := vault.NewClient(profile.Address, token)
	if err != nil {
		return nil, fmt.Errorf("vault client: %w", err)
	}
	return New(client, profile.EnvFile, profile.Paths, profile.Backup), nil
}

// Run fetches all secrets and merges them into the env file.
// If backup is enabled a timestamped copy is created before writing.
func (s *Syncer) Run() (*env.DiffResult, error) {
	if s.backup {
		if _, err := env.Backup(s.envPath); err != nil {
			return nil, fmt.Errorf("backup: %w", err)
		}
	}

	existing, err := env.Read(s.envPath)
	if err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}

	incoming := make(map[string]string)
	for _, p := range s.paths {
		secrets, err := s.client.GetSecrets(p)
		if err != nil {
			return nil, fmt.Errorf("get secrets %s: %w", p, err)
		}
		for k, v := range secrets {
			incoming[k] = v
		}
	}

	diffResult := env.Diff(existing, incoming)

	w := env.NewWriter(s.envPath)
	if err := w.Merge(incoming); err != nil {
		return nil, fmt.Errorf("write env: %w", err)
	}

	return &diffResult, nil
}
