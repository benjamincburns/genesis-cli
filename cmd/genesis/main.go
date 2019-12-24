package main

import(
	"log"

	"github.com/whiteblock/genesis-cli/pkg/oauth2-noserver"

	"golang.org/x/oauth2"
)
func main(){
	conf := &oauth2.Config{
		ClientID:     "cli",
		Scopes:       []string{"offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://auth.infra.whiteblock.io/auth/realms/wb/protocol/openid-connect/auth",
			TokenURL: "https://auth.infra.whiteblock.io/auth/realms/wb/protocol/openid-connect/token",
		},
	}

	client, err := oauth2ns.AuthenticateUser(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("token %s\n",*client.Token)
}