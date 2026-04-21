package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/env"
)

var swapCmd = &cobra.Command{
	Use:   "swap --file <env-file> OLD=NEW [OLD2=NEW2 ...]",
	Short: "Rename keys in a local .env file",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSwap,
}

var (
	swapFile          string
	swapFailOnMissing bool
	swapDryRun        bool
)

func init() {
	swapCmd.Flags().StringVar(&swapFile, "file", ".env", "path to the .env file")
	swapCmd.Flags().BoolVar(&swapFailOnMissing, "fail-on-missing", false, "return error when source key is absent")
	swapCmd.Flags().BoolVar(&swapDryRun, "dry-run", false, "print result without writing changes")
	rootCmd.AddCommand(swapCmd)
}

func runSwap(cmd *cobra.Command, args []string) error {
	secrets, err := env.Read(swapFile)
	if err != nil {
		return fmt.Errorf("swap: read %s: %w", swapFile, err)
	}

	pairs := make(map[string]string, len(args))
	for _, arg := range args {
		for i := 0; i < len(arg); i++ {
			if arg[i] == '=' {
				pairs[arg[:i]] = arg[i+1:]
				break
			}
		}
	}
	if len(pairs) == 0 {
		return fmt.Errorf("swap: no valid OLD=NEW pairs provided")
	}

	out, res, err := env.Swap(secrets, env.SwapOptions{
		Pairs:         pairs,
		FailOnMissing: swapFailOnMissing,
		DryRun:        swapDryRun,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), res.Summary())

	if swapDryRun {
		return nil
	}

	w := env.NewWriter(swapFile)
	if err := w.Write(out); err != nil {
		return fmt.Errorf("swap: write %s: %w", swapFile, err)
	}

	if len(res.Missing) > 0 {
		for _, k := range res.Missing {
			fmt.Fprintf(os.Stderr, "swap: warning: key %q not found\n", k)
		}
	}
	return nil
}
