package auth

import (
	"fmt"
	"log/slog"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func newAzureOauth2Config(AuthPort int, TenantID string, ClientID string) *oauth2.Config {
	redirectURL := fmt.Sprintf("http://localhost:%d/oauth/callback", AuthPort)
	slog.Debug("redirectURL", "value", redirectURL, "method", "newAzureOauth2Config")
	scopes := []string{fmt.Sprintf("%s/.default", ClientID)}
	slog.Debug("scopes", "value", scopes, "method", "newAzureOauth2Config")
	return &oauth2.Config{
		ClientID:    ClientID,
		Endpoint:    microsoft.AzureADEndpoint(TenantID),
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}
}
