package auth

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lcrownover/hpcadmin-lib/pkg/oauth"
)

func (h *AuthHandler) LoadAccessToken() (string, bool) {
	// load a local access token
	slog.Debug("loading local access token", "method", "LoadAccessToken")
	t, err := h.readAccessTokenFromFile()
	if err != nil {
		slog.Debug(fmt.Sprintf("error reading local access token: %v", err), "method", "LoadAccessToken")
		return "", false
	}
	jwtToken, err := oauth.GetJWTFromTokenString(t)
	if err != nil {
		slog.Debug(fmt.Sprintf("error getting JWT from token: %v", err), "method", "LoadAccessToken")
		return "", false
	}
	if !oauth.JWTTokenIsValid(jwtToken) {
		slog.Debug("token is expired or invalid", "method", "LoadAccessToken")
		return "", false
	}
	slog.Debug("local access token loaded", "method", "LoadAccessToken")
	return t, true
}

func (h *AuthHandler) readAccessTokenFromFile() (string, error) {
	// load a local access token
	slog.Debug("reading local access token from file", "method", "readAccessTokenFromFile")
	t, err := os.ReadFile(h.ConfigDir + "/token")
	if err != nil {
		slog.Debug("error reading local access token from file", "method", "readAccessTokenFromFile")
		return "", err
	}
	slog.Debug("local access token read from file", "method", "readAccessTokenFromFile")
	return string(t), nil
}

func (h *AuthHandler) SaveAccessToken(token string) error {
	// save the token to a file
	slog.Debug("saving access token", "method", "SaveAccessToken")
	err := os.WriteFile(h.ConfigDir+"/token", []byte(token), 0600)
	if err != nil {
		slog.Debug("error saving access token", "method", "SaveAccessToken")
		return err
	}
	slog.Debug("access token saved", "method", "SaveAccessToken")
	return nil
}
