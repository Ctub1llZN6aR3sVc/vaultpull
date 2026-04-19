package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/env"
)

var intersectKeep string

var intersectCmd = &cobra.Command{
	Use:   "intersect <file-a> <file-b>",
	Short: "Output only keys present in both env files",
	Args:  cobra.ExactArgs(2),
	RunE:  runIntersect,
}

func init() {
	intersectCmd.Flags().StringVar(&intersectKeep, "keep", "a", "Which file's value to keep when key exists in both (a or b)")
	rootCmd.AddCommand(intersectCmd)
}

func runIntersect(cmd *cobra.Command, args []string) error {
	a, err := env.Read(args[0])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}

	b, err := env.Read(args[1])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[1], err)
	}

	opts := env.IntersectOptions{KeepValue: intersectKeep}
	res := env.Intersect(a, b, opts)

	fmt.Println(res.Summary())
	for _, k := range res.Kept {
		fmt.Printf("  %s=%s\n", k, res.Secrets[k])
	}

	return nil
}
