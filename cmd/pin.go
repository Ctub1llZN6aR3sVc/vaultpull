package cmd

import (
	"fmt"
	"time"

	"github.com/densestvoid/vaultpull/internal/env"
	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Manage pinned secret values",
}

var pinAddCmd = &cobra.Command{
	Use:   "add <key> <value>",
	Short: "Pin a key to a specific value",
	Args:  cobra.ExactArgs(2),
	RunE:  runPinAdd,
}

var pinListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pinned keys",
	RunE:  runPinList,
}

var pinFile string
var pinTTL time.Duration

func init() {
	pinAddCmd.Flags().StringVar(&pinFile, "pin-file", "pins.json", "path to pin file")
	pinAddCmd.Flags().DurationVar(&pinTTL, "ttl", 0, "optional TTL for the pin (e.g. 24h)")
	pinListCmd.Flags().StringVar(&pinFile, "pin-file", "pins.json", "path to pin file")
	pinCmd.AddCommand(pinAddCmd)
	pinCmd.AddCommand(pinListCmd)
	rootCmd.AddCommand(pinCmd)
}

func runPinAdd(cmd *cobra.Command, args []string) error {
	pins, err := env.LoadPins(pinFile)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	entry := env.PinEntry{Key: args[0], Value: args[1], PinnedAt: now}
	if pinTTL > 0 {
		entry.ExpiresAt = now.Add(pinTTL)
	}
	pins = append(pins, entry)
	if err := env.SavePins(pinFile, pins); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Pinned %s\n", args[0])
	return nil
}

func runPinList(cmd *cobra.Command, args []string) error {
	pins, err := env.LoadPins(pinFile)
	if err != nil {
		return err
	}
	if len(pins) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No pinned keys.")
		return nil
	}
	for _, p := range pins {
		expiry := "no expiry"
		if !p.ExpiresAt.IsZero() {
			expiry = p.ExpiresAt.Format(time.RFC3339)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s  (pinned: %s, expires: %s)\n",
			p.Key, p.PinnedAt.Format(time.RFC3339), expiry)
	}
	return nil
}
