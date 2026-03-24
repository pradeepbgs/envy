package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pradeepbgs/envy/internal/config"
	"github.com/pradeepbgs/envy/internal/storage"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete an env set from R2",
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config not found — run 'envy init' first: %w", err)
	}

	fmt.Printf("Delete '%s' from R2? [y/N]: ", name)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(answer)) != "y" {
		fmt.Println("Aborted.")
		return nil
	}

	r2 := storage.New(cfg)
	if err := r2.Delete(cmd.Context(), name+".enc"); err != nil {
		return err
	}

	fmt.Printf("Deleted '%s'\n", name)
	return nil
}
