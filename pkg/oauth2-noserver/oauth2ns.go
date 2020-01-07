/**
 * Based on original which can be found at the url: https://github.com/nmrshll/oauth2-noserver
 */

package oauth2ns

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/message"
	"github.com/whiteblock/genesis-cli/pkg/util"

	rndm "github.com/nmrshll/rndm-go"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

type AuthorizedClient struct {
	*http.Client
	Token *oauth2.Token
}

var (
	gconf            = config.NewConfig()
	ErrAuthTimeout   = fmt.Errorf(message.AuthenticationTimedOut)
	ErrNilAuthConfig = stacktrace.NewError("oauthConfig can't be nil")
)

const (

	// seconds to wait before giving up on auth and exiting
	oauthStateStringContextKey = 987
)

type AuthenticateUserOption func(*AuthenticateUserFuncConfig) error
type AuthenticateUserFuncConfig struct {
	AuthCallHTTPParams url.Values
}

func WithAuthCallHTTPParams(values url.Values) AuthenticateUserOption {
	return func(conf *AuthenticateUserFuncConfig) error {
		conf.AuthCallHTTPParams = values
		return nil
	}
}

// AuthenticateUser starts the login process
func AuthenticateUser(oauthConfig *oauth2.Config, options ...AuthenticateUserOption) (*AuthorizedClient, error) {
	// validate params
	if oauthConfig == nil {
		return nil, ErrNilAuthConfig
	}
	// read options
	var optionsConfig AuthenticateUserFuncConfig
	for _, processConfigFunc := range options {
		processConfigFunc(&optionsConfig)
	}

	// add transport for self-signed certificate to context
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, sslcli)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	oauthConfig.RedirectURL = fmt.Sprintf("http://%s", gconf.RedirectURL)

	// Some random string, random for each request
	oauthStateString := rndm.String(8)
	ctx = context.WithValue(ctx, oauthStateStringContextKey, oauthStateString)
	urlString := oauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)

	if optionsConfig.AuthCallHTTPParams != nil {
		parsedURL, err := url.Parse(urlString)
		if err != nil {
			return nil, stacktrace.Propagate(err, "failed parsing url string")
		}
		params := parsedURL.Query()
		for key, value := range optionsConfig.AuthCallHTTPParams {
			params[key] = value
		}
		parsedURL.RawQuery = params.Encode()
		urlString = parsedURL.String()
	}

	clientChan, stopHTTPServerChan, cancelAuthentication := startHTTPServer(ctx, oauthConfig)
	util.AuthPrint(message.PreBrowserAuthShowURL)
	util.AuthPrint("\"" + urlString + "\"")
	time.Sleep(1000 * time.Millisecond)
	err := open.Run(urlString)
	if err != nil {
		util.ErrorFatal(message.FailedToOpenBrowserTab)
	}
	time.Sleep(600 * time.Millisecond)

	// shutdown the server after timeout
	go func() {
		time.Sleep(gconf.AuthTimeout)
		stopHTTPServerChan <- struct{}{}
	}()

	select {
	// wait for client on clientChan
	case client := <-clientChan:
		// After the callbackHandler returns a client, it's time to shutdown the server gracefully
		stopHTTPServerChan <- struct{}{}
		return client, nil

		// if authentication process is cancelled first return an error
	case <-cancelAuthentication:
		return nil, ErrAuthTimeout
	}
}

func startHTTPServer(ctx context.Context, conf *oauth2.Config) (clientChan chan *AuthorizedClient,
	stopHTTPServerChan chan struct{}, cancelAuthentication chan struct{}) {
	// init returns
	clientChan = make(chan *AuthorizedClient)
	stopHTTPServerChan = make(chan struct{})
	cancelAuthentication = make(chan struct{})

	http.HandleFunc("/", callbackHandler(ctx, conf, clientChan))
	srv := &http.Server{Addr: gconf.RedirectURL}

	// handle server shutdown signal
	go func() {
		// wait for signal on stopHTTPServerChan
		<-stopHTTPServerChan
		log.Trace("Shutting down server...")

		// give it 5 sec to shutdown gracefully, else quit program
		d := time.Now().Add(5 * time.Second)
		ctx, cancel := context.WithDeadline(context.Background(), d)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			util.Errorf("auth server could not shutdown gracefully: %v", err)
		}

		// after server is shutdown, quit program
		cancelAuthentication <- struct{}{}
	}()

	// handle callback request
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			util.ErrorFatalf("listen: %s\n", err)
		}
		log.Debug("Server gracefully stopped")
	}()

	return clientChan, stopHTTPServerChan, cancelAuthentication
}

func callbackHandler(ctx context.Context, oauthConfig *oauth2.Config,
	clientChan chan *AuthorizedClient) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		requestStateString := ctx.Value(oauthStateStringContextKey).(string)
		responseStateString := r.FormValue("state")
		if responseStateString != requestStateString {
			util.Errorf("invalid oauth state, expected '%s', got '%s'\n",
				requestStateString, responseStateString)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		code := r.FormValue("code")
		token, err := oauthConfig.Exchange(ctx, code)
		if err != nil {
			util.Errorf("oauthoauthConfig.Exchange() failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// The HTTP Client returned by oauthConfig.Client will refresh the token as necessary
		client := &AuthorizedClient{
			oauthConfig.Client(ctx, token),
			token,
		}
		// show success page
		successPage := `
		<div></div>
		<script>window.onload=function(){setTimeout(this.close, 60)}</script>
		`
		fmt.Fprintf(w, successPage)
		// quitSignalChan <- quitSignal
		clientChan <- client
	}
}
