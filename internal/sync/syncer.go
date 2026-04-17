package sync

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

// VaultClient fetches secrets from a Vault path.
type VaultClient interface {
	GetSecrets(path string) (map[string]string, error)
}

// Syncer pulls secrets from Vault and writes them to an env file.
type Syncer struct {
	client  VaultClient
	paths   []string
	output  string

	Backup    *env.BackupOptions
	Audit     string
	Encrypt   []byte
	Snapshot  string
	Lint      bool
	Validate  bool
	Template  bool
	Filter    *env.FilterOptions
	Transform *env.TransformOptions
	Resolve   *env.ResolveOptions
	PinFile   string
	Expiry    map[string]string
}

// New creates a Syncer with the given client, paths, and output file.
func New(client VaultClient, paths []string, output string) *Syncer {
	return &Syncer{client: client, paths: paths, output: output}
}

// NewFromConfig builds a Syncer from a config profile.
func NewFromConfig(client VaultClient, profile config.Profile) *Syncer {
	s := New(client, profile.Paths, profile.Output)
	if profile.Backup.Enabled {
		s.Backup = &env.BackupOptions{
			Dir:     profile.Backup.Dir,
			MaxKeep: profile.Backup.MaxKeep,
		}
	}
	s.Audit = profile.AuditLog
	s.Snapshot = profile.Snapshot
	s.Lint = profile.Lint
	s.Validate = profile.Validate
	s.Template = profile.Template
	return s
}

// Run executes the sync: fetch, resolve, transform, filter, write.
func (s *Syncer) Run() error {
	merged := make(map[string]string)

	for _, path := range s.paths {
		secrets, err := s.client.GetSecrets(path)
		if err != nil {
			return fmt.Errorf("fetch %q: %w", path, err)
		}
		for k, v := range secrets {
			merged[k] = v
		}
	}

	// Apply pins.
	if s.PinFile != "" {
		pins, err := env.LoadPins(s.PinFile)
		if err != nil {
			return err
		}
		if pins != nil {
			merged = env.ApplyPins(merged, pins)
		}
	}

	// Resolve defaults / required.
	if s.Resolve != nil {
		result, err := env.Resolve(merged, *s.Resolve)
		if err != nil {
			return err
		}
		merged = result.Secrets
	}

	// Template rendering.
	if s.Template {
		var err error
		merged, err = env.RenderMap(merged)
		if err != nil {
			return err
		}
	}

	// Filter.
	if s.Filter != nil {
		merged = env.Filter(merged, *s.Filter)
	}

	// Transform.
	if s.Transform != nil {
		merged = env.Transform(merged, *s.Transform)
	}

	// Validate.
	if s.Validate {
		if _, err := env.Validate(merged); err != nil {
			return err
		}
	}

	// Lint.
	if s.Lint {
		if _, err := env.Lint(merged); err != nil {
			return err
		}
	}

	// Backup existing file.
	if s.Backup != nil {
		if err := env.Backup(s.output, *s.Backup); err != nil {
			return err
		}
	}

	// Encrypt.
	if s.Encrypt != nil {
		var err error
		merged, err = env.EncryptSecrets(merged, s.Encrypt)
		if err != nil {
			return err
		}
	}

	// Read existing, compute diff, write.
	existing, _ := env.Read(s.output)
	diff := env.Diff(existing, merged)
	_ = diff

	w := env.NewWriter(s.output)
	if err := w.Write(merged); err != nil {
		return err
	}

	// Snapshot.
	if s.Snapshot != "" {
		if err := env.SaveSnapshot(s.Snapshot, merged); err != nil {
			return err
		}
	}

	// Audit.
	if s.Audit != "" {
		entry := env.AuditEntry{Keys: keysOf(merged)}
		if err := env.WriteAuditLog(s.Audit, entry); err != nil {
			return err
		}
	}

	return nil
}

func keysOf(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
