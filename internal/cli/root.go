package cli

import (
	"fmt"
	"log/slog"

	"github.com/lcrownover/hpcadmin-cli/internal/config"
	"github.com/lcrownover/hpcadmin-cli/internal/util"
	"github.com/spf13/cobra"
)

var (
	debug bool

	rootCmd = &cobra.Command{
		Use:   "hpcadmin",
		Short: "HPCAdmin CLI",
		Long:  `HPCAdmin is a CMDB for hosting membership information for the Talapas HPC at University of Oregon`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			util.ConfigureLogging(debug)
			slog.Debug("Starting hpcadmin-cli", "method", "Execute")

			_, err := config.GetCLIConfig()
			if err != nil {
				util.PrintAndExit(fmt.Sprintf("Error getting CLI config: %v\n", err), 1)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
}

func Execute() error {
	return rootCmd.Execute()
}
