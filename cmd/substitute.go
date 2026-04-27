package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/env"
)

var substituteCmd = &cobra.Command{
	Use:   "substitute <file>",
	Short: "Expand ${VAR} references within a .env file using values from the same file",
	Args:  cobra.ExactArgs(1),
	RunE:  runSubstitute,
}

func init() {
	substituteCmd.Flags().Bool("strict", false, "Fail if a referenced variable is not found")
	substituteCmd.Flags().String("fallback", "", "Default value to use for unresolved references")
	substituteCmd.Flags().Bool("dry-run", false, "Print result without writing to file")
	rootCmd.AddCommand(substituteCmd)
}

func runSubstitute(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	strict, _ := cmd.Flags().GetBool("strict")
	fallback, _ := cmd.Flags().GetString("fallback")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	secrets, err := env.Read(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	opts := &env.SubstituteOptions{
		Strict:   strict,
		Fallback: fallback,
		DryRun:   dryRun,
	}

	out, result, err := env.Substitute(secrets, opts)
	if err != nil {
		return fmt.Errorf("substitution failed: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), result.Summary())

	if dryRun {
		for _, k := range result.Substituted {
			fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s = %s\n", k, out[k])
		}
		return nil
	}

	w := env.NewWriter(filePath)
	if err := w.Write(out); err != nil {
		return fmt.Errorf("writing %s: %w", filePath, err)
	}

	if len(result.Unresolved) > 0 {
		fmt.Fprintf(os.Stderr, "warning: unresolved variables: %v\n", result.Unresolved)
	}
	return nil
}
