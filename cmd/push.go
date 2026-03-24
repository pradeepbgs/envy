package cmd

import (
	"fmt"
	"os"

	"github.com/pradeepbgs/envy/internal/config"
	"github.com/pradeepbgs/envy/internal/crypto"
	"github.com/pradeepbgs/envy/internal/storage"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push <name> <envfile>",
	Short: "Encrypt and upload a .env file to R2",
	Args:  cobra.ExactArgs(2),
	RunE:  runPush,
}

func runPush(cmd *cobra.Command, args []string) error {
	name := args[0]
	envFile := args[1]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config not found — run 'envy init' first: %w", err)
	}

	data, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", envFile, err)
	}

	key, err := config.DecodeKey(cfg.EncryptionKey)
	if err != nil {
		return fmt.Errorf("invalid encryption key in config: %w", err)
	}

	encrypted, err := crypto.Encrypt(data, key)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	r2 := storage.New(cfg)
	if err := r2.Upload(cmd.Context(), name+".enc", encrypted); err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	fmt.Printf("Pushed '%s' to R2\n", name)
	return nil
}
