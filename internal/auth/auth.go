package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/browser"

	"github.com/lcrownover/hpcadmin-cli/internal/util"
	"golang.org/x/oauth2"
)

type AuthVendor int

const (
	Azure  AuthVendor = iota
	Google AuthVendor = iota
)

type OauthHandlerOptions struct {
	Vendor              AuthVendor
	TenantID            string
	ClientID            string
	SkipTLSVerification bool
}

func NewOauthHandlerOptions(vendor AuthVendor, tenantID string, clientID string) OauthHandlerOptions {
	return OauthHandlerOptions{
		Vendor:   vendor,
		TenantID: tenantID,
		ClientID: clientID,
	}
}

type AuthHandler struct {
	Ctx          context.Context
	ConfigDir    string
	ListenAddr   string
	Oauth2Config *oauth2.Config
	HttpClient   *http.Client
	HttpServer   *http.Server
	HttpMux      *http.ServeMux
	AccessToken  string
	AuthDoneCh   chan struct{}
}

func getRandomPort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port
}

func NewAuthHandler(configDirectory string, opts OauthHandlerOptions) *AuthHandler {
	slog.Debug("creating new auth handler", "method", "NewAuthHandler")
	slog.Debug("configDirectory", "value", configDirectory, "method", "NewAuthHandler")
	ctx := context.Background()
	authPort := getRandomPort()
	slog.Debug("authPort", "value", authPort, "method", "NewAuthHandler")

	// oauth2 config includes things like the client id,
	// the auth endpoint, redirectURL, and scopes
	var oauthConfig *oauth2.Config
	switch opts.Vendor {
	case Azure:
		oauthConfig = newAzureOauth2Config(authPort, opts.TenantID, opts.ClientID)
	default:
		util.PrintAndExit(fmt.Sprintf("Error: %v is not a valid auth vendor\n", opts.Vendor), 1)
	}

	// register a custom http client that maybe skips SSL verification
	// and store it in ctx
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: opts.SkipTLSVerification},
	}
	sslclient := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslclient)

	// create a new http server and mux
	mux := http.NewServeMux()
	server := &http.Server{Addr: fmt.Sprintf(":%d", authPort), Handler: mux}
	return &AuthHandler{
		Ctx:          ctx,
		ConfigDir:    configDirectory,
		Oauth2Config: oauthConfig,
		HttpClient:   sslclient,
		HttpMux:      mux,
		HttpServer:   server,
		AccessToken:  "",
		AuthDoneCh:   make(chan struct{}, 1),
	}
}

func (h *AuthHandler) GetAuthenticationURL() string {
	return h.Oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (h *AuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("callbackHandler called", "method", "CallbackHandler")
	slog.Debug("parsing query string", "method", "CallbackHandler")
	queryParts, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		slog.Error(err.Error())
	}
	code := queryParts["code"][0]

	// Exchange will do the handshake to retrieve the initial access token.
	slog.Debug("exchanging code for token", "method", "CallbackHandler")
	tok, err := h.Oauth2Config.Exchange(h.Ctx, code)
	if err != nil {
		slog.Debug(err.Error())
		fmt.Fprintf(w, "Authentication Code exchange failed")
		os.Exit(1)
	}
	h.AccessToken = tok.AccessToken

	// The HTTP Client returned by conf.Client will refresh the token as necessary.
	client := h.Oauth2Config.Client(h.Ctx, tok)
	h.HttpClient = client

	// show succes page
	slog.Debug("showing success page", "method", "CallbackHandler")
	successHTML := `
<h1>Authentication Success</h1>
<p>You are authenticated and can now return to the CLI.</p>
`
	fmt.Fprint(w, successHTML)
	slog.Debug("sending auth done signal", "method", "CallbackHandler")
	h.AuthDoneCh <- struct{}{}
	slog.Debug("callbackHandler finished", "method", "CallbackHandler")
}

func (h *AuthHandler) Authenticate() (string, error) {
	var err error
	util.InfoPrint("You will now be taken to your browser for authentication")

	time.Sleep(1 * time.Second)

	url := h.GetAuthenticationURL()
	err = browser.OpenURL(url)
	if err != nil {
		return "", fmt.Errorf("error opening browser: %v", err)
	}

	time.Sleep(1 * time.Second)

	go func() {
		h.HttpMux.HandleFunc("/oauth/callback", h.CallbackHandler)
		slog.Debug("starting server", "method", "Authenticate")
		err := h.HttpServer.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				// ErrServerClosed is returned by ListenAndServe and is fine
				slog.Error("server error", "method", "Authenticate")
			}
		}
	}()

	for n := 0; n < 1; {
		select {
		case <-h.AuthDoneCh:
			slog.Debug("authentication successful, shutting down server", "method", "Authenticate")
			h.HttpServer.Shutdown(h.Ctx)
			slog.Debug("server shut down", "method", "Authenticate")
			n++
		case <-time.After(1 * time.Minute):
			slog.Debug("authentication failed, shutting down server", "method", "Authenticate")
			h.HttpServer.Shutdown(h.Ctx)
			slog.Debug("server shut down", "method", "Authenticate")
			return "", fmt.Errorf("authentication timed out")
		}
	}

	return h.AccessToken, nil
}
