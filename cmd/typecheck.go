package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/env"
)

var typecheckFile string

func init() {
	typecheckCmd := &cobra.Command{
		Use:   "typecheck",
		Short: "Validate .env values against expected types",
		RunE:  runTypeCheck,
	}
	typecheckCmd.Flags().StringVarP(&typecheckFile, "file", "f", ".env", "path to .env file")
	rootCmd.AddCommand(typecheckCmd)
}

func runTypeCheck(cmd *cobra.Command, args []string) error {
	secrets, err := env.Read(typecheckFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", typecheckFile, err)
	}

	if len(secrets) == 0 {
		fmt.Println("typecheck: no secrets found")
		return nil
	}

	// Build rules from args: KEY=TYPE pairs
	var rules []env.TypeRule
	for _, arg := range args {
		var key, typ string
		if n, _ := fmt.Sscanf(arg, "%s", &key); n != 1 {
			continue
		}
		for i, c := range arg {
			if c == '=' {
				key = arg[:i]
				typ = arg[i+1:]
				break
			}
		}
		if key != "" && typ != "" {
			rules = append(rules, env.TypeRule{Key: key, Expected: typ})
		}
	}

	if len(rules) == 0 {
		fmt.Println("typecheck: no rules provided (pass KEY=TYPE args)")
		return nil
	}

	res := env.TypeCheck(secrets, rules)
	fmt.Println(res.Summary())
	if !res.OK() {
		os.Exit(1)
	}
	return nil
}
