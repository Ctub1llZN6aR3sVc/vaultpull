package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint secrets fetched from Vault for common issues",
	RunE:  runLint,
}

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().String("config", "vaultpull.yaml", "config file path")
	lintCmd.Flags().String("profile", "", "profile name")
	lintCmd.Flags().String("token", "", "Vault token (overrides VAULT_TOKEN)")
}

func runLint(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	profileName, _ := cmd.Flags().GetString("profile")
	token, _ := cmd.Flags().GetString("token")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("get profile: %w", err)
	}
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("vault token required (--token or VAULT_TOKEN)")
	}

	client, err := vault.NewClient(profile.Address, token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	merged := map[string]string{}
	for _, p := range profile.Paths {
		secrets, err := client.GetSecrets(p)
		if err != nil {
			return fmt.Errorf("fetch %s: %w", p, err)
		}
		for k, v := range secrets {
			merged[k] = v
		}
	}

	result := env.Lint(merged)
	if result.IsClean() {
		fmt.Println("✔", result.Summary())
		return nil
	}
	for _, issue := range result.Issues {
		fmt.Fprintf(os.Stderr, "  ✘ %s\n", issue)
	}
	return fmt.Errorf("%s", result.Summary())
}
