package cmd

import (
	"fmt"
	"strings"

	"github.com/pradeepbgs/envy/internal/config"
	"github.com/pradeepbgs/envy/internal/storage"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all env sets stored in R2",
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config not found — run 'envy init' first: %w", err)
	}

	r2 := storage.New(cfg)
	keys, err := r2.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("list failed: %w", err)
	}

	if len(keys) == 0 {
		fmt.Println("No env sets found.")
		return nil
	}

	fmt.Println("Available env sets:")
	for _, k := range keys {
		fmt.Printf("  %s\n", strings.TrimSuffix(k, ".enc"))
	}
	return nil
}
