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

var (
	conf         = config.NewConfig()
	globalClient *oauth2ns.AuthorizedClient
)

func GetToken() *oauth2.Token {
	token := new(oauth2.Token)

	if len(conf.GenesisCredentials) != 0 {
		err := json.Unmarshal([]byte(conf.GenesisCredentials), token)
		if err == nil {
			return token
		}
	}
	if !conf.UserDir.Exists(conf.TokenFile) {
		log.Trace("the token file does not exist")
		return nil
	}
	data, err := conf.UserDir.ReadFile(conf.TokenFile)
	if err != nil {
		log.WithField("error", err).Debug("couldn't read token file")
		return nil
	}

	err = json.Unmarshal(data, token)
	if err != nil {
		log.WithField("error", err).Debug("couldn't parse token file")
		return nil
	}
	return token
}

func getClientFromLocalToken(authConf *oauth2.Config) *oauth2ns.AuthorizedClient {
	token := GetToken()
	if token == nil {
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

func Login() (*oauth2ns.AuthorizedClient, error) {
	client, err := oauth2ns.AuthenticateUser(getAuthConf())
	if err != nil {
		return nil, err
	}
	err = storeToken(client)
	if err != nil {
		log.WithField("error", err).Error("failed to store token")
	}
	// Check that the user exists
	err = Get(conf.GetSelfURL(), nil)
	if err == nil {
		return client, nil
	}
	globalClient = client
	// User does not exist
	util.Print("Setting up your account...")
	_, err = Post(conf.CreateUserURL(), []byte(``))
	if err != nil {
		util.Errorf("first time setup failed: %s", err.Error())
	}
	return client, err
}

func GetClient() (*oauth2ns.AuthorizedClient, error) {
	if globalClient != nil {
		return globalClient, nil
	}
	authConf := getAuthConf()

	var err error

	globalClient = getClientFromLocalToken(authConf)
	if globalClient != nil {
		return globalClient, nil
	}

	globalClient, err = Login()
	if err != nil {
		return nil, err
	}
	return globalClient, nil

}

func getAuthConf() *oauth2.Config {
	return &oauth2.Config{
		ClientID: "cli",
		Scopes:   []string{"offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s%s", conf.AuthEndpoint, conf.AuthPath),
			TokenURL: fmt.Sprintf("https://%s%s", conf.AuthEndpoint, conf.TokenPath),
		},
	}
}
