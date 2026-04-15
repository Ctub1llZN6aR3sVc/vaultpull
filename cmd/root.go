package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	profile string
	envFile string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpull pulls secrets from HashiCorp Vault and writes them
into local .env files with support for multiple profiles.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "vaultpull.yaml", "config file path")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "profile to use")
	rootCmd.PersistentFlags().StringVarP(&envFile, "env-file", "e", ".env", "target .env file path")
}
