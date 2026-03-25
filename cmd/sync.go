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

// runSync's JOB is to load the config , take user input about name , target dir then download the encrypted file from r2 , decrypt that and save in the given dir

func runSync(cmd *cobra.Command, args []string) error {
	name := args[0]
	targetdir := args[1]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config not found — run 'envy init' first: %w", err)
	}

	// so here we will check if the given dir already has .env file
	outPath := filepath.Join(targetdir, ".env")
	if !force {
		if _, err := os.Stat(outPath); err == nil {
			return fmt.Errorf("%s already exists — use --force to overwrite", outPath)
		}
	}

	// decode the key
	key, err := config.Decodekey(cfg.Encryptionkey)
	if err != nil {
		return fmt.Errorf("invalid encryption key in config: %w", err)
	}

	// get the R2
	r2 := storage.New(cfg)

	//downalod the encrypted file by given name
	encrypted, err := r2.Download(cmd.Context(), name+".enc")
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// decyrpt the file
	plainText, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		return fmt.Errorf("decrypt failed: %w", err)
	}

	// write to the target dir
	if err := os.MkdirAll(targetdir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}
	if err := os.WriteFile(outPath, plainText, 0600); err != nil {
		return fmt.Errorf("failed to write .env: %w", err)
	}

	fmt.Printf("Synced '%s' -> %s\n", name, outPath)
	return nil
}
