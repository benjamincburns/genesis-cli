package auth

import (
	"fmt"

	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/oauth2-noserver"

	"golang.org/x/oauth2"
)

var conf = config.NewConfig()

func GetClient() (*oauth2ns.AuthorizedClient, error) {
	authConf := &oauth2.Config{
		ClientID: "cli",
		Scopes:   []string{"offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s%s", conf.AuthEndpoint, conf.AuthPath),
			TokenURL: fmt.Sprintf("https://%s%s", conf.AuthEndpoint, conf.TokenPath),
		},
	}

	return oauth2ns.AuthenticateUser(authConf)
}
