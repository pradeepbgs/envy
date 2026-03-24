package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pradeepbgs/envy/internal/config"
	"github.com/pradeepbgs/envy/internal/crypto"
	"github.com/pradeepbgs/envy/internal/storage"
	"github.com/spf13/cobra"
)

var force bool

var syncCmd = &cobra.Command{
	Use:   "sync <name> <targetdir>",
	Short: "Download and decrypt a .env file from R2 into a directory",
	Args:  cobra.ExactArgs(2),
	RunE:  runSync,
}

func init() {
	syncCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing .env")
}

func runSync(cmd *cobra.Command, args []string) error {
	name := args[0]
	targetDir := args[1]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config not found — run 'envy init' first: %w", err)
	}

	outPath := filepath.Join(targetDir, ".env")
	if !force {
		if _, err := os.Stat(outPath); err == nil {
			return fmt.Errorf("%s already exists — use --force to overwrite", outPath)
		}
	}

	r2 := storage.New(cfg)
	encrypted, err := r2.Download(cmd.Context(), name+".enc")
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	key, err := config.DecodeKey(cfg.EncryptionKey)
	if err != nil {
		return fmt.Errorf("invalid encryption key in config: %w", err)
	}

	plaintext, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(outPath, plaintext, 0600); err != nil {
		return fmt.Errorf("failed to write .env: %w", err)
	}

	fmt.Printf("Synced '%s' -> %s\n", name, outPath)
	return nil
}
