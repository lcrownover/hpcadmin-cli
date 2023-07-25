/*
Copyright © 2023 Lucas Crownover <lcrownover127@gmail.com>
*/
package cmd

import (
	"github.com/lcrownover/hpcadmin-cli/internal/core"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a user",
	Long: `Get details for a user.`,
	Run: func(cmd *cobra.Command, args []string) {
		core.CLIUserShow(cmd.Flag("username").Value.String(), cmd.Flag("email").Value.String(), cmd.Flag("firstname").Value.String(), cmd.Flag("lastname").Value.String())
	},
}

func init() {
	userCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	showCmd.Flags().StringP("username", "u", "", "Username")
	showCmd.MarkFlagRequired("username")
	showCmd.Flags().StringP("email", "e", "", "Email")
	showCmd.MarkFlagRequired("email")
	showCmd.Flags().StringP("firstname", "f", "", "First Name")
	showCmd.MarkFlagRequired("firstname")
	showCmd.Flags().StringP("lastname", "l", "", "Last Name")
	showCmd.MarkFlagRequired("lastname")
}
