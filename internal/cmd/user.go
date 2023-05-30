/*
Copyright © 2023 Lucas Crownover <lcrownover127@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	RunE: requireSubcommand,
	Use:   "user",
	Short: "Manage HPC users",
	Long: `User objects are the basis of HPC management.

There are three roles a user can serve for a PIRG:
1. Owner: Responsible for the project, almost always the PI.
2. Admin: Capable of the day-to-day management of the PIRG.
3. Member: Member of the PIRG, and has access to the PIRG's resources.`,
}

func init() {
	rootCmd.AddCommand(userCmd)
}
