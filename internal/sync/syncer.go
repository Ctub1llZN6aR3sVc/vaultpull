package sync

import (
	"fmt"

	"github.com/elizaos/vaultpull/internal/config"
	"github.com/elizaos/vaultpull/internal/env"
	"github.com/elizaos/vaultpull/internal/vault"
)

// VaultClient abstracts secret retrieval.
type VaultClient interface {
	GetSecrets(path string) (map[string]string, error)
}

// Options configures a Syncer run.
type Options struct {
	Paths          []string
	OutputFile     string
	Merge          bool
	BackupEnabled  bool
	BackupDir      string
	BackupKeep     int
	AuditLog       string
	Validate       bool
	ExpiryTTLDays  int
	RenderTemplate bool
	Filter         env.FilterOptions
	Transform      env.TransformOptions
	EncryptKey     []byte
}

// Syncer pulls secrets from Vault and writes them to an env file.
type Syncer struct {
	client  VaultClient
	options Options
}

// New creates a Syncer with the given client and options.
func New(client VaultClient, opts Options) *Syncer {
	return &Syncer{client: client, options: opts}
}

// NewFromConfig builds a Syncer from a loaded config profile.
func NewFromConfig(cfg *config.Config, profile string) (*Syncer, error) {
	p, err := cfg.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	client, err := vault.NewClient(p.Address, p.Token)
	if err != nil {
		return nil, err
	}
	opts := Options{
		Paths:      p.Paths,
		OutputFile: p.OutputFile,
		Merge:      p.Merge,
	}
	return New(client, opts), nil
}

// Run executes the sync: fetch, transform, filter, encrypt, write.
func (s *Syncer) Run() error {
	merged := make(map[string]string)
	for _, path := range s.options.Paths {
		secrets, err := s.client.GetSecrets(path)
		if err != nil {
			return fmt.Errorf("path %s: %w", path, err)
		}
		for k, v := range secrets {
			merged[k] = v
		}
	}

	if s.options.RenderTemplate {
		var err error
		merged, err = env.RenderMap(merged)
		if err != nil {
			return err
		}
	}

	merged = env.Filter(merged, s.options.Filter)
	merged = env.Transform(merged, s.options.Transform)

	if s.options.Validate {
		result := env.Validate(merged)
		if !result.Valid {
			return fmt.Errorf("validation failed: %s", result.Summary)
		}
	}

	if len(s.options.EncryptKey) > 0 {
		var err error
		merged, err = env.EncryptSecrets(merged, s.options.EncryptKey)
		if err != nil {
			return fmt.Errorf("encrypt secrets: %w", err)
		}
	}

	if s.options.BackupEnabled {
		if err := env.Rotate(s.options.OutputFile, env.RotateOptions{
			BackupDir:  s.options.BackupDir,
			KeepBackups: s.options.BackupKeep,
		}); err != nil {
			return err
		}
	}

	w := env.NewWriter(s.options.OutputFile)
	var writeErr error
	if s.options.Merge {
		writeErr = w.Merge(merged)
	} else {
		writeErr = w.Write(merged)
	}
	if writeErr != nil {
		return writeErr
	}

	if s.options.AuditLog != "" {
		_ = env.WriteAuditLog(s.options.AuditLog, merged)
	}

	return nil
}
