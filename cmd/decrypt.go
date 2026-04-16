package cmd

import (
	"fmt"
	"os"

	"github.com/elizaos/vaultpull/internal/env"
	"github.com/spf13/cobra"
)

var decryptKeyHex string

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt an encrypted .env file in place",
	RunE:  runDecrypt,
}

func init() {
	decryptCmd.Flags().StringVar(&decryptKeyHex, "key", "", "32-byte AES key (hex or raw)")
	decryptCmd.Flags().StringP("file", "f", ".env", "path to encrypted .env file")
	_ = decryptCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(decryptCmd)
}

func runDecrypt(cmd *cobra.Command, _ []string) error {
	file, _ := cmd.Flags().GetString("file")
	key := []byte(decryptKeyHex)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return fmt.Errorf("key must be 16, 24, or 32 bytes; got %d", len(key))
	}

	secrets, err := env.Read(file)
	if err != nil {
		return fmt.Errorf("read %s: %w", file, err)
	}
	if len(secrets) == 0 {
		fmt.Fprintln(os.Stderr, "no secrets found")
		return nil
	}

	decrypted, err := env.DecryptSecrets(secrets, key)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	w := env.NewWriter(file)
	if err := w.Write(decrypted); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	fmt.Printf("decrypted %d keys into %s\n", len(decrypted), file)
	return nil
}
