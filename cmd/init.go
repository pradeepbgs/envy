package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pradeepbgs/envy/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize envy with your R2 credentials",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("R2 Endpoint (https://<account_id>.r2.cloudflarestorage.com): ")
	endpoint, _ := reader.ReadString('\n')
	endpoint = strings.TrimSpace(endpoint)

	fmt.Print("Access Key ID: ")
	accessKey, _ := reader.ReadString('\n')
	accessKey = strings.TrimSpace(accessKey)

	fmt.Print("Secret Access Key: ")
	secretKey, _ := reader.ReadString('\n')
	secretKey = strings.TrimSpace(secretKey)

	fmt.Print("Bucket name [envy-store]: ")
	bucket, _ := reader.ReadString('\n')
	bucket = strings.TrimSpace(bucket)
	if bucket == "" {
		bucket = "envy-store"
	}

	encKey, err := config.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	cfg := &config.Config{
		R2Endpoint:    endpoint,
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		Bucket:        bucket,
		EncryptionKey: encKey,
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Config saved to ~/.envy/config.json")
	fmt.Println("Encryption key generated and stored. Keep ~/.envy/config.json safe — losing it means losing access to your envs.")
	return nil
}
