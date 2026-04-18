package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"vaultpull/internal/vault"
)

// groupFlagKey is the flag name for group rules.
const groupFlagKey = "group"

// registerGroupFlags attaches the --group flag to a command.
func registerGroupFlags(cmd *cobra.Command) {
	cmd.Flags().StringArray(groupFlagKey, nil, "Group secrets by pattern: pattern=group (e.g. DB_*=database)")
}

// buildGroupRules reads and parses --group flags from the command.
func buildGroupRules(cmd *cobra.Command) ([]vault.GroupRule, error) {
	raws, err := cmd.Flags().GetStringArray(groupFlagKey)
	if err != nil || len(raws) == 0 {
		return nil, err
	}
	rules, err := vault.ParseGroupRules(raws)
	if err != nil {
		return nil, fmt.Errorf("--group: %w", err)
	}
	return rules, nil
}
