package cmd

import (
	"fmt"
	"os"

	"github.com/eliziario/vaultpull/internal/env"
	"github.com/spf13/cobra"
)

var aliasFile string
var aliasRemoveSource bool
var aliasDryRun bool
var aliasMappings []string

func init() {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Copy secret values to alias keys in a .env file",
		RunE:  runAlias,
	}
	aliasCmd.Flags().StringVarP(&aliasFile, "file", "f", ".env", "Target .env file")
	aliasCmd.Flags().StringArrayVarP(&aliasMappings, "map", "m", nil, "Alias mapping as NEW_KEY=SRC_KEY (repeatable)")
	aliasCmd.Flags().BoolVar(&aliasRemoveSource, "remove-source", false, "Remove source key after aliasing")
	aliasCmd.Flags().BoolVar(&aliasDryRun, "dry-run", false, "Preview changes without writing")
	rootCmd.AddCommand(aliasCmd)
}

func runAlias(cmd *cobra.Command, args []string) error {
	if len(aliasMappings) == 0 {
		return fmt.Errorf("at least one --map NEW_KEY=SRC_KEY is required")
	}

	secrets, err := env.Read(aliasFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", aliasFile, err)
	}

	aliases := make(map[string]string, len(aliasMappings))
	for _, m := range aliasMappings {
		for i, c := range m {
			if c == '=' {
				aliases[m[:i]] = m[i+1:]
				break
			}
		}
	}

	out, res := env.Alias(secrets, env.AliasOptions{
		Aliases:      aliases,
		RemoveSource: aliasRemoveSource,
		DryRun:       aliasDryRun,
	})

	fmt.Fprintln(cmd.OutOrStdout(), res.Summary())
	if len(res.Missing) > 0 {
		for _, k := range res.Missing {
			fmt.Fprintf(cmd.OutOrStdout(), "  missing source key: %s\n", k)
		}
	}

	if aliasDryRun {
		return nil
	}

	w := env.NewWriter(aliasFile)
	if err := w.Write(out); err != nil {
		return fmt.Errorf("write %s: %w", aliasFile, err)
	}
	_ = out
	_ = os.Stdout
	return nil
}
