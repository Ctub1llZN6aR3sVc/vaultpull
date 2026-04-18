package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/env/reader"
)

var castCmd = &cobra.Command{
	Use:   "cast [file]",
	Short: "Cast env values to specified types and print result",
	Args:  cobra.ExactArgs(1),
	RunE:  runCast,
}

var castTypes []string
var castStrict bool

func init() {
	castCmd.Flags().StringArrayVar(&castTypes, "type", nil, "KEY=TYPE pairs (e.g. PORT=int)")
	castCmd.Flags().BoolVar(&castStrict, "strict", false, "fail on uncastable value")
	rootCmd.AddCommand(castCmd)
}

func runCast(cmd *cobra.Command, args []string) error {
	file := args[0]
	secrets, err := env.Read(file)
	if err != nil {
		return fmt.Errorf("read %s: %w", file, err)
	}

	types := make(map[string]string)
	for _, pair := range castTypes {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --type value %q, expected KEY=TYPE", pair)
		}
		types[parts[0]] = parts[1]
	}

	out, results, err := env.Cast(secrets, env.CastOptions{Types: types, Strict: castStrict})
	if err != nil {
		return err
	}

	for _, r := range results {
		fmt.Fprintf(os.Stderr, "cast: %s %q -> %q (%s)\n", r.Key, r.Original, r.Casted, r.Type)
	}

	for k, v := range out {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
	}
	_ = reader.Noop()
	return nil
}
