package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

var tagFilter []string

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "List secrets filtered by tag",
		RunE:  runTag,
	}
	tagCmd.Flags().StringVar(&cfgFile, "config", "vaultpull.yaml", "config file")
	tagCmd.Flags().StringVar(&profile, "profile", "", "profile name")
	tagCmd.Flags().StringArrayVar(&tagFilter, "tag", nil, "required tag (repeatable)")
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	p, err := cfg.GetProfile(profile)
	if err != nil {
		return err
	}
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return fmt.Errorf("VAULT_TOKEN is not set")
	}
	client, err := vault.NewClient(p.Address, token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	secrets := make(map[string]string)
	for _, path := range p.Paths {
		data, err := client.GetSecrets(path)
		if err != nil {
			return fmt.Errorf("get secrets %s: %w", path, err)
		}
		for k, v := range data {
			secrets[k] = v
		}
	}
	tagMap := p.Tags
	if tagMap == nil {
		tagMap = map[string]string{}
	}
	var filtered map[string]string
	if len(tagFilter) > 0 {
		filtered = env.FilterByTag(secrets, tagMap, tagFilter...)
	} else {
		filtered = secrets
	}
	result := env.Tag(filtered, tagMap)
	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", env.SummaryByTag(result))
	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		tags := result.Tagged[k]
		if len(tags) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s (untagged)\n", k)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s [%s]\n", k, joinTags(tags))
		}
	}
	return nil
}

func joinTags(tags []string) string {
	out := ""
	for i, t := range tags {
		if i > 0 {
			out += ", "
		}
		out += t
	}
	return out
}
