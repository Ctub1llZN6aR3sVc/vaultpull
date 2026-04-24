package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/env"
)

var reorderCmd = &cobra.Command{
	Use:   "reorder",
	Short: "Reorder keys in a .env file",
	RunE:  runReorder,
}

func init() {
	reorderCmd.Flags().String("file", ".env", "path to .env file")
	reorderCmd.Flags().StringSlice("keys", nil, "explicit key order (comma-separated)")
	reorderCmd.Flags().Bool("alpha", false, "sort keys alphabetically")
	reorderCmd.Flags().Bool("dry-run", false, "print result without writing")
	rootCmd.AddCommand(reorderCmd)
}

func runReorder(cmd *cobra.Command, _ []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	alpha, _ := cmd.Flags().GetBool("alpha")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	secrets, err := env.Read(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}

	opts := env.ReorderOptions{
		Keys:         keys,
		Alphabetical: alpha,
		DryRun:       dryRun,
	}

	out, result := env.Reorder(secrets, opts)

	if dryRun {
		for _, k := range result.Ordered {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, out[k])
		}
		fmt.Fprintln(cmd.OutOrStdout(), "# "+result.Summary())
		return nil
	}

	w := env.NewWriter(filePath)
	if err := w.Write(out); err != nil {
		return fmt.Errorf("write %s: %w", filePath, err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), result.Summary())

	if len(result.Ordered) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "order: %s\n", strings.Join(result.Ordered, ", "))
	}

	_ = os
	return nil
}
