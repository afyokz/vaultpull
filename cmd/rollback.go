package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/vault"
)

var (
	rollbackMaxSize int
	rollbackLabel   string
	rollbackID      string
)

var rollbackStore = vault.NewRollbackStore(vault.DefaultRollbackMaxSize)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Manage secret rollback points",
}

var rollbackListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available rollback entries",
	RunE:  runRollbackList,
}

var rollbackApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a rollback entry by ID (prints secrets to stdout)",
	RunE:  runRollbackApply,
}

func init() {
	rollbackCmd.AddCommand(rollbackListCmd)
	rollbackCmd.AddCommand(rollbackApplyCmd)
	rollbackApplyCmd.Flags().StringVar(&rollbackID, "id", "", "Rollback entry ID to apply")
	_ = rollbackApplyCmd.MarkFlagRequired("id")
	RootCmd.AddCommand(rollbackCmd)
}

func runRollbackList(cmd *cobra.Command, args []string) error {
	entries := rollbackStore.List()
	fmt.Print(vault.FormatRollbackList(entries))
	return nil
}

func runRollbackApply(cmd *cobra.Command, args []string) error {
	entry, ok := rollbackStore.Get(rollbackID)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: rollback entry %q not found\n", rollbackID)
		os.Exit(1)
	}
	fmt.Printf("# Rollback entry: %s (%s)\n", entry.ID, entry.Timestamp.Format("2006-01-02T15:04:05Z"))
	for k, v := range entry.Secrets {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
