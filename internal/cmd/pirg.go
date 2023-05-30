/*
Copyright © 2023 Lucas Crownover <lcrownover127@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// pirgCmd represents the pirg command
var pirgCmd = &cobra.Command{
	RunE:  requireSubcommand,
	Use:   "pirg",
	Short: "Manage HPC PIRGs",
	Long: `PIRGs (PI Research Groups) represent a group of assigned users and resources.

There are three roles a user can serve for a PIRG:
1. Owner: Responsible for the project, almost always the PI.
2. Admin: Capable of the day-to-day management of the PIRG.
3. Member: Member of the PIRG, and has access to the PIRG's resources.`,
}

func init() {
	rootCmd.AddCommand(pirgCmd)
}
