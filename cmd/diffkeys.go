package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

func init() {
	diffKeysCmd := &cobra.Command{
		Use:   "diffkeys <file-a> <file-b>",
		Short: "Compare key sets between two .env files",
		Args:  cobra.ExactArgs(2),
		RunE:  runDiffKeys,
	}
	rootCmd.AddCommand(diffKeysCmd)
}

func runDiffKeys(cmd *cobra.Command, args []string) error {
	a, err := env.Read(args[0])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}
	b, err := env.Read(args[1])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[1], err)
	}

	d := env.DiffKeys(a, b)

	w := cmd.OutOrStdout()
	if len(d.OnlyInA) == 0 && len(d.OnlyInB) == 0 {
		fmt.Fprintln(w, "Key sets are identical.")
		return nil
	}

	if len(d.OnlyInA) > 0 {
		fmt.Fprintf(w, "Only in %s:\n", args[0])
		for _, k := range d.OnlyInA {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}
	if len(d.OnlyInB) > 0 {
		fmt.Fprintf(w, "Only in %s:\n", args[1])
		for _, k := range d.OnlyInB {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}
	if len(d.InBoth) > 0 {
		fmt.Fprintf(w, "Shared (%d): %s\n", len(d.InBoth), strings.Join(d.InBoth, ", "))
	}

	fmt.Fprintln(w, env.KeyDiffSummary(d))
	_ = os.Stderr
	return nil
}
