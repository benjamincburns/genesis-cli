package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/oauth2-noserver"
	"github.com/whiteblock/genesis-cli/pkg/util"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var conf = config.NewConfig()

func getClientFromLocalToken(authConf *oauth2.Config) *oauth2ns.AuthorizedClient {
	if !conf.UserDir.Exists(conf.TokenFile) {
		return nil
	}
	data, err := conf.UserDir.ReadFile(conf.TokenFile)
	if err != nil {
		log.WithField("error", err).Debug("couldn't read token file")
		return nil
	}

	token := new(oauth2.Token)

	err = json.Unmarshal(data, token)
	if err != nil {
		log.WithField("error", err).Debug("couldn't parse token file")
		return nil
	}

	return &oauth2ns.AuthorizedClient{
		Client: authConf.Client(context.Background(), token),
		Token:  token,
	}

}

func storeToken(client *oauth2ns.AuthorizedClient) error {
	data, err := json.Marshal(client.Token)
	if err != nil {
		return err
	}
	return conf.UserDir.WriteFile(conf.TokenFile, data)
}

func GetClient() (*oauth2ns.AuthorizedClient, error) {
	authConf := &oauth2.Config{
		ClientID: "cli",
		Scopes:   []string{"offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s%s", conf.AuthEndpoint, conf.AuthPath),
			TokenURL: fmt.Sprintf("https://%s%s", conf.AuthEndpoint, conf.TokenPath),
		},
	}

	client := getClientFromLocalToken(authConf)
	if client != nil {
		return client, nil
	}

	client, err := oauth2ns.AuthenticateUser(authConf)
	if err != nil {
		return nil, err
	}
	err = storeToken(client)
	if err != nil {
		util.Errorf("couldn't store token: %v", err)
	}
	return client, nil

}
