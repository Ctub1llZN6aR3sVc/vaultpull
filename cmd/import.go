package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/env"
)

var importOverwrite bool
var importDryRun bool

func init() {
	importCmd := &cobra.Command{
		Use:   "import [file]",
		Short: "Import secrets from a local .env file into the active profile's env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runImport,
	}
	importCmd.Flags().StringVar(&cfgFile, "config", ".vaultpull.yaml", "config file")
	importCmd.Flags().StringVar(&profileName, "profile", "", "profile to use")
	importCmd.Flags().BoolVar(&importOverwrite, "overwrite", false, "overwrite existing keys")
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "preview changes without writing")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	src := args[0]

	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("get profile: %w", err)
	}

	dst, err := env.Read(profile.Output)
	if err != nil {
		return fmt.Errorf("read output file: %w", err)
	}

	opts := env.ImportOptions{
		Overwrite: importOverwrite,
		DryRun:    importDryRun,
	}
	res, err := env.ImportFromFile(src, dst, opts)
	if err != nil {
		return err
	}

	if !importDryRun {
		w := env.NewWriter(profile.Output)
		if err := w.Write(dst); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}

	sort.Strings(res.Keys)
	for _, k := range res.Keys {
		if importDryRun {
			fmt.Printf("[dry-run] would import: %s\n", k)
		} else {
			fmt.Printf("imported: %s\n", k)
		}
	}
	fmt.Printf("done: %d imported, %d skipped\n", res.Imported, res.Skipped)
	return nil
}
