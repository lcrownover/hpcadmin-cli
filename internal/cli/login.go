package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/lcrownover/hpcadmin-cli/internal/auth"
	"github.com/lcrownover/hpcadmin-cli/internal/config"
	"github.com/lcrownover/hpcadmin-cli/internal/util"
	"github.com/spf13/cobra"
)

type AuthInfo struct {
	TenantID string `json:"tenant_id"`
	ClientID string `json:"client_id"`
}

func init() {
	rootCmd.AddCommand(LoginCmd)
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to HPCAdmin",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("running login command", "method", "LoginCmd.Run")
		cfg, err := config.GetCLIConfig()
		if err != nil {
			util.PrintAndExit(fmt.Sprintf("Error getting CLI config: %v\n", err), 1)
		}
		
		authInfo, err := GetAuthInfo(cfg)
		if err != nil {
			util.ErrorPrint(fmt.Sprintf("Error getting auth info: %v\n", err))
			os.Exit(1)
		}

		var accessToken string
		azureAuthOptions := auth.NewOauthHandlerOptions(auth.Azure, authInfo.TenantID, authInfo.ClientID)

		ah := auth.NewAuthHandler(cfg.ConfigDir, azureAuthOptions)
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

func GetAuthInfo(config *config.CLIConfig) (*AuthInfo, error) {
	var authInfo AuthInfo
	slog.Debug("getting auth info from hpcadmin server", "method", "GetAuthInfo")
	resp, err := http.Get(fmt.Sprintf("%s/oauth/info", config.BaseURL))
	if err != nil {
		slog.Debug("failed to get auth info from hpcadmin server", "method", "GetAuthInfo", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	slog.Debug("successfully retrieved auth info from hpcadmin server", "method", "GetAuthInfo")
	slog.Debug("reading response body", "method", "GetAuthInfo")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Debug("error reading response body", "method", "GetAuthInfo", "error", err)
		return nil, err
	}
	slog.Debug("successfully read response body", "method", "GetAuthInfo")
	slog.Debug("unmarshalling response body", "method", "GetAuthInfo")
	err = json.Unmarshal(body, &authInfo)
	if err != nil {
		slog.Debug("error unmarshalling response body", "method", "GetAuthInfo", "error", err)
		return nil, err
	}
	slog.Debug("successfully retrieved auth info from hpcadmin server", "info", fmt.Sprintf("%+v", &authInfo), "method", "GetAuthInfo")
	return &authInfo, nil
}
