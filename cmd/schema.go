package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/sync"
)

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate synced secrets against a schema",
		RunE:  runSchema,
	}
	schemaCmd.Flags().StringP("config", "c", "vaultpull.yaml", "config file")
	schemaCmd.Flags().StringP("profile", "p", "", "profile name")
	schemaCmd.Flags().StringP("token", "t", "", "vault token (or set VAULT_TOKEN)")
	schemaCmd.Flags().Bool("strict", false, "fail on unexpected keys")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	profileName, _ := cmd.Flags().GetString("profile")
	token, _ := cmd.Flags().GetString("token")
	strict, _ := cmd.Flags().GetBool("strict")

	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("vault token required via --token or VAULT_TOKEN")
	}

	syncer, err := sync.NewFromConfig(cfgPath, profileName, token)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	secrets, err := syncer.FetchSecrets()
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	fields := syncer.SchemaFields()
	if len(fields) == 0 {
		fmt.Println("no schema defined in profile, skipping")
		return nil
	}

	result := env.ValidateSchema(secrets, fields, strict)
	fmt.Println(result.Summary())
	if result.HasError {
		return fmt.Errorf("schema validation failed")
	}
	return nil
}
