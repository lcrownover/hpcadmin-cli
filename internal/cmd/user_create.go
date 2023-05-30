/*
Copyright © 2023 Lucas Crownover <lcrownover127@gmail.com>
*/
package cmd

import (
	"github.com/lcrownover/hpcadmin-cli/internal/core"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a user",
	Long: `Users must have a valid University account and email address.
`,
	Run: func(cmd *cobra.Command, args []string) {
		core.CLIUserCreate(cmd.Flag("username").Value.String(), cmd.Flag("email").Value.String(), cmd.Flag("firstname").Value.String(), cmd.Flag("lastname").Value.String())
	},
}

func init() {
	userCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringP("username", "u", "", "Username")
	createCmd.MarkFlagRequired("username")
	createCmd.Flags().StringP("email", "e", "", "Email")
	createCmd.MarkFlagRequired("email")
	createCmd.Flags().StringP("firstname", "f", "", "First Name")
	createCmd.MarkFlagRequired("firstname")
	createCmd.Flags().StringP("lastname", "l", "", "Last Name")
	createCmd.MarkFlagRequired("lastname")
}
