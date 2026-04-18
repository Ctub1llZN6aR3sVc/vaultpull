package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/env"
)

var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Resolve secrets by chaining multiple .env files (first wins)",
	RunE:  runChain,
}

func init() {
	chainCmd.Flags().StringSlice("files", nil, "Ordered list of .env files to chain (highest priority first)")
	chainCmd.Flags().Bool("show-source", false, "Annotate each key with its source file")
	rootCmd.AddCommand(chainCmd)
}

func runChain(cmd *cobra.Command, _ []string) error {
	files, _ := cmd.Flags().GetStringSlice("files")
	showSource, _ := cmd.Flags().GetBool("show-source")

	if len(files) == 0 {
		return fmt.Errorf("at least one --files entry is required")
	}

	sources := make([]map[string]string, 0, len(files))
	for _, f := range files {
		m, err := env.Read(f)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("reading %s: %w", f, err)
		}
		sources = append(sources, m)
	}

	merged := env.ChainAll(sources...)

	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if showSource {
			src := sourceOf(k, files, sources)
			fmt.Printf("%s=%s  # %s\n", k, merged[k], src)
		} else {
			fmt.Printf("%s=%s\n", k, merged[k])
		}
	}
	return nil
}

func sourceOf(key string, files []string, sources []map[string]string) string {
	for i, src := range sources {
		if v, ok := src[key]; ok && v != "" {
			return files[i]
		}
	}
	return "unknown"
}
