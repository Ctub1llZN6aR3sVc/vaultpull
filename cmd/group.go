package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/env"
)

var groupPrefixes []string
var groupSep string
var groupFile string

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Partition .env keys by prefix groups",
		RunE:  runGroup,
	}
	groupCmd.Flags().StringSliceVar(&groupPrefixes, "prefix", nil, "Prefixes to group by (repeatable)")
	groupCmd.Flags().StringVar(&groupSep, "sep", "_", "Separator between prefix and key")
	groupCmd.Flags().StringVar(&groupFile, "file", ".env", "Source .env file")
	rootCmd.AddCommand(groupCmd)
}

func runGroup(cmd *cobra.Command, args []string) error {
	if len(groupPrefixes) == 0 {
		return fmt.Errorf("at least one --prefix is required")
	}

	secrets, err := env.Read(groupFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", groupFile, err)
	}

	res := env.Group(secrets, groupPrefixes, groupSep)

	groupNames := make([]string, 0, len(res.Groups))
	for k := range res.Groups {
		groupNames = append(groupNames, k)
	}
	sort.Strings(groupNames)

	for _, g := range groupNames {
		fmt.Fprintf(os.Stdout, "[%s]\n", g)
		keys := make([]string, 0, len(res.Groups[g]))
		for k := range res.Groups[g] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(os.Stdout, "  %s=%s\n", k, res.Groups[g][k])
		}
	}

	if len(res.Ungrouped) > 0 {
		fmt.Fprintln(os.Stdout, "[ungrouped]")
		ukeys := make([]string, 0, len(res.Ungrouped))
		for k := range res.Ungrouped {
			ukeys = append(ukeys, k)
		}
		sort.Strings(ukeys)
		for _, k := range ukeys {
			fmt.Fprintf(os.Stdout, "  %s=%s\n", k, res.Ungrouped[k])
		}
	}

	_ = strings.Join(groupPrefixes, ",") // suppress unused import
	return nil
}
