package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

var watchInterval int

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch <env-file>",
		Short: "Poll an env file and report when it changes",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatch,
	}
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 5, "poll interval in seconds")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	path := args[0]

	w, err := env.NewWatchState(path)
	if err != nil {
		return fmt.Errorf("watch: init: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Watching %s (interval: %ds)\n", path, watchInterval)

	ticker := time.NewTicker(time.Duration(watchInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			changed, err := w.Changed()
			if err != nil {
				return fmt.Errorf("watch: check: %w", err)
			}
			if changed {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s has changed\n",
					time.Now().Format(time.RFC3339), path)
				if err := w.Refresh(); err != nil {
					return fmt.Errorf("watch: refresh: %w", err)
				}
			}
		case <-cmd.Context().Done():
			return nil
		}
	}
}
