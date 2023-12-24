package cli

import (
	"fmt"
	"os"

	"github.com/lcrownover/hpcadmin-cli/internal/auth"
	"github.com/lcrownover/hpcadmin-cli/internal/util"
	"github.com/spf13/cobra"
)

// TODO(lcrown): move these to a config file
const AZURE_TENANT_ID = "8f0b198f-f447-4cfe-ba03-526b46c661f8"
const AZURE_CLIENT_ID = "1951f213-c370-4a77-b7cd-7a4c303df45a"

func init() {
	rootCmd.AddCommand(LoginCmd)
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to HPCAdmin",
	Run: func(cmd *cobra.Command, args []string) {
		var accessToken string
		azureAuthOptions := auth.NewOauthHandlerOptions(auth.Azure, AZURE_TENANT_ID, AZURE_CLIENT_ID)

		ah := auth.NewAuthHandler(configDir, azureAuthOptions)
		accessToken, ok := ah.LoadAccessToken()
		if !ok {
			accessToken, err = ah.Authenticate()
			if err != nil {
				util.ErrorPrint(fmt.Sprintf("Error authenticating: %v\n", err))
				os.Exit(1)
			}
			ah.SaveAccessToken(accessToken)
		}

		// token can be accessed
		fmt.Printf("token: %v\n", accessToken)
	},
}
