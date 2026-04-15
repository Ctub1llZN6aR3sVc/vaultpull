package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Profile represents a named Vault sync configuration.
type Profile struct {
	Name      string `yaml:"name"`
	VaultAddr string `yaml:"vault_addr"`
	VaultPath string `yaml:"vault_path"`
	OutputFile string `yaml:"output_file"`
	AuthMethod string `yaml:"auth_method"` // token, approle
}

// Config holds all vaultpull configuration.
type Config struct {
	DefaultProfile string             `yaml:"default_profile"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

// Load reads and parses the config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found at %q: run 'vaultpull init' to create one", path)
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}

	return &cfg, nil
}

// GetProfile returns the profile by name, or the default profile if name is empty.
func (c *Config) GetProfile(name string) (Profile, error) {
	key := name
	if key == "" {
		key = c.DefaultProfile
	}
	if key == "" {
		return Profile{}, fmt.Errorf("no profile specified and no default_profile set")
	}
	p, ok := c.Profiles[key]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found in config", key)
	}
	return p, nil
}
