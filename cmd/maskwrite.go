package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/env"
)

var maskwriteCmd = &cobra.Command{
	Use:   "maskwrite <file>",
	Short: "Mask sensitive values in a .env file in place",
	Args:  cobra.ExactArgs(1),
	RunE:  runMaskWrite,
}

var (
	maskwriteAutoDetect bool
	maskwriteKeys       []string
	maskwritePlaceholder string
	maskwriteDryRun     bool
)

func init() {
	maskwriteCmd.Flags().BoolVar(&maskwriteAutoDetect, "auto", true, "Auto-detect sensitive keys")
	maskwriteCmd.Flags().StringSliceVar(&maskwriteKeys, "keys", nil, "Explicit keys to mask")
	maskwriteCmd.Flags().StringVar(&maskwritePlaceholder, "placeholder", "***", "Placeholder for masked values")
	maskwriteCmd.Flags().BoolVar(&maskwriteDryRun, "dry-run", false, "Print result without writing")
	rootCmd.AddCommand(maskwriteCmd)
}

func runMaskWrite(cmd *cobra.Command, args []string) error {
	file := args[0]

	secrets, err := env.Read(file)
	if err != nil {
		return fmt.Errorf("read %s: %w", file, err)
	}

	opts := env.MaskWriteOptions{
		AutoDetect:  maskwriteAutoDetect,
		Keys:        maskwriteKeys,
		Placeholder: maskwritePlaceholder,
		DryRun:      maskwriteDryRun,
	}

	out, result := env.MaskWrite(secrets, opts)
	fmt.Fprintln(cmd.ErrOrStderr(), result.Summary())

	if maskwriteDryRun {
		for k, v := range out {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
		}
		return nil
	}

	w := env.NewWriter(file)
	if err := w.Write(out); err != nil {
		return fmt.Errorf("write %s: %w", file, err)
	}

	fmt.Fprintf(os.Stdout, "wrote masked secrets to %s\n", file)
	return nil
}
