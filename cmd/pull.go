package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/config"
	"github.com/vaultpull/vaultpull/internal/sync"
	"github.com/vaultpull/vaultpull/internal/vault"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull secrets from Vault and write to .env file",
	Long:  `Connects to HashiCorp Vault using the specified profile and syncs secrets into the target .env file.`,
	RunE:  runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	p, err := cfg.GetProfile(profile)
	if err != nil {
		return fmt.Errorf("profile %q not found: %w", profile, err)
	}

	token := os.Getenv("VAULT_TOKEN")
	if p.Token != "" {
		token = p.Token
	}
	if token == "" {
		return fmt.Errorf("vault token not set: use VAULT_TOKEN env var or set token in profile")
	}

	client, err := vault.NewClient(p.Address, token)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	outFile := envFile
	if p.EnvFile != "" {
		outFile = p.EnvFile
	}

	syncer := sync.NewFromConfig(client, outFile, p.Paths)
	if err := syncer.Run(); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	fmt.Printf("Successfully synced secrets to %s\n", outFile)
	return nil
}
