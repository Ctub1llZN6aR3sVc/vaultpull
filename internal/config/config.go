package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Profile holds the settings for a single named environment.
type Profile struct {
	Address string   `yaml:"address"`
	EnvFile string   `yaml:"env_file"`
	Paths   []string `yaml:"paths"`
	Backup  bool     `yaml:"backup"`
}

// Config is the top-level configuration structure.
type Config struct {
	DefaultProfile string             `yaml:"default_profile"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}

// GetProfile returns the named profile, falling back to DefaultProfile when
// name is empty. Returns an error if the profile is not found.
func (c *Config) GetProfile(name string) (*Profile, error) {
	if name == "" {
		name = c.DefaultProfile
	}
	p, ok := c.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", name)
	}
	return &p, nil
}
