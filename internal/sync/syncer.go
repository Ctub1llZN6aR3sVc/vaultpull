package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

// SecretClient abstracts Vault secret fetching.
type SecretClient interface {
	GetSecrets(path string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets and writing env files.
type Syncer struct {
	client  SecretClient
	profile config.Profile
	output  io.Writer
}

// New creates a Syncer with the given client and profile.
func New(client SecretClient, profile config.Profile, output io.Writer) *Syncer {
	if output == nil {
		output = os.Stdout
	}
	return &Syncer{client: client, profile: profile, output: output}
}

// NewFromConfig builds a Syncer from a loaded config and profile name.
func NewFromConfig(cfg *config.Config, profileName string, output io.Writer) (*Syncer, error) {
	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return nil, err
	}
	client, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		return nil, err
	}
	return New(client, profile, output), nil
}

// Run fetches secrets for all paths and merges them into the env file.
func (s *Syncer) Run() (*env.DiffResult, error) {
	before, _ := env.Read(s.profile.EnvFile)

	merged := make(map[string]string)
	for k, v := range before {
		merged[k] = v
	}

	for _, path := range s.profile.Paths {
		secrets, err := s.client.GetSecrets(path)
		if err != nil {
			return nil, fmt.Errorf("fetching path %q: %w", path, err)
		}
		for k, v := range secrets {
			merged[k] = v
		}
	}

	w := env.NewWriter(s.profile.EnvFile)
	if err := w.Write(merged); err != nil {
		return nil, fmt.Errorf("writing env file: %w", err)
	}

	diff := env.Diff(before, merged)
	env.PrintDiff(s.output, diff)
	return diff, nil
}
