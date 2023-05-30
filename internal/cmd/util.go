package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func requireSubcommand(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("%s requires a subcommand", cmd.Name())
}
