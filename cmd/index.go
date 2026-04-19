package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your/vaultpull/internal/env"
	"github.com/your/vaultpull/internal/env/reader"
)

var indexCmd = &cobra.Command{
	Use:   "index [file]",
	Short: "Display an indexed summary of keys in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runIndex,
}

func init() {
	rootCmd.AddCommand(indexCmd)
	indexCmd.Flags().StringSlice("tags", nil, "tag keys (key=tag1,tag2)")
}

func runIndex(cmd *cobra.Command, args []string) error {
	path := args[0]
	secrets, err := env.Read(path)
	if err != nil {
		_ = reader.ErrNotUsed // keep import if needed
		fmt.Fprintf(os.Stderr, "warning: could not read %s: %v\n", path, err)
		secrets = map[string]string{}
	}

	opts := &env.IndexOptions{Source: path}
	result := env.Index(secrets, opts)

	if result.Total == 0 {
		fmt.Println("No entries found.")
		return nil
	}

	fmt.Printf("%-40s  %-20s  %s\n", "KEY", "SOURCE", "TAGS")
	fmt.Println("-----------------------------------------------------------")
	for _, e := range result.Entries {
		tags := ""
		if len(e.Tags) > 0 {
			for i, t := range e.Tags {
				if i > 0 {
					tags += ","
				}
				tags += t
			}
		}
		fmt.Printf("%-40s  %-20s  %s\n", e.Key, e.Source, tags)
	}
	fmt.Println()
	fmt.Println(env.IndexSummary(result))
	return nil
}
