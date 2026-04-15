package sync

import (
	"fmt"

	"github.com/eliziario/vaultpull/internal/config"
	"github.com/eliziario/vaultpull/internal/env"
	"github.com/eliziario/vaultpull/internal/vault"
)

// VaultClient defines the interface for fetching secrets.
type VaultClient interface {
	GetSecrets(path string) (map[string]string, error)
}

// Options configures optional syncer behaviour.
type Options struct {
	DryRun          bool
	BackupEnabled   bool
	BackupDir       string
	BackupKeep      int
	AuditLogPath    string
	Validate        bool
	TTL             string
	RenderTemplates bool
}

// Syncer pulls secrets from Vault and writes them to an env file.
type Syncer struct {
	client  VaultClient
	envFile string
	paths   []string
	opts    Options
}

// New creates a Syncer with explicit dependencies.
func New(client VaultClient, envFile string, paths []string, opts Options) *Syncer {
	return &Syncer{client: client, envFile: envFile, paths: paths, opts: opts}
}

// NewFromConfig constructs a Syncer from a loaded config profile.
func NewFromConfig(client VaultClient, profile *config.Profile) *Syncer {
	return &Syncer{
		client:  client,
		envFile: profile.EnvFile,
		paths:   profile.Paths,
		opts: Options{
			BackupEnabled:   profile.Backup.Enabled,
			BackupDir:       profile.Backup.Dir,
			BackupKeep:      profile.Backup.Keep,
			AuditLogPath:    profile.AuditLog,
			Validate:        profile.Validate,
			TTL:             profile.TTL,
			RenderTemplates: profile.RenderTemplates,
		},
	}
}

// Run fetches secrets from all configured paths and writes them to the env file.
func (s *Syncer) Run() error {
	merged := make(map[string]string)

	for _, path := range s.paths {
		secrets, err := s.client.GetSecrets(path)
		if err != nil {
			return fmt.Errorf("vault path %q: %w", path, err)
		}
		for k, v := range secrets {
			merged[k] = v
		}
	}

	if s.opts.RenderTemplates {
		rendered, err := env.RenderMap(merged)
		if err != nil {
			return fmt.Errorf("template rendering: %w", err)
		}
		merged = rendered
	}

	if s.opts.Validate {
		result := env.Validate(merged)
		if !result.Valid {
			return fmt.Errorf("validation failed: %s", result.Summary)
		}
	}

	if s.opts.BackupEnabled {
		if err := env.Rotate(s.envFile, env.RotateOptions{
			BackupDir:  s.opts.BackupDir,
			KeepBackups: s.opts.BackupKeep,
		}); err != nil {
			return fmt.Errorf("backup: %w", err)
		}
	}

	w := env.NewWriter(s.envFile)
	if err := w.Merge(merged); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}

	if s.opts.AuditLogPath != "" {
		keys := make([]string, 0, len(merged))
		for k := range merged {
			keys = append(keys, k)
		}
		_ = env.WriteAuditLog(s.opts.AuditLogPath, vault.TokenSource(), s.paths, keys)
	}

	return nil
}
