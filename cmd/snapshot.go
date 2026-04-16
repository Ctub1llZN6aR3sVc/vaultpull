package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/env"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Show diff between a saved snapshot and the current .env file",
	RunE:  runSnapshot,
}

func init() {
	snapshotCmd.Flags().StringP("snapshot", "s", ".vaultpull.snap.json", "Path to snapshot file")
	snapshotCmd.Flags().StringP("env", "e", ".env", "Path to .env file to compare")
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, _ []string) error {
	snapPath, _ := cmd.Flags().GetString("snapshot")
	envPath, _ := cmd.Flags().GetString("env")

	snap, err := env.LoadSnapshot(snapPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	current, err := env.Read(envPath)
	if err != nil {
		return fmt.Errorf("read env: %w", err)
	}

	diff := env.DiffSnapshot(snap, current)
	env.PrintDiff(cmd.OutOrStdout(), diff)
	return nil
}
