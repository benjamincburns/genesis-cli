package main

import(
	"log"

	"github.com/whiteblock/genesis-cli/pkg/oauth2-noserver"

	"golang.org/x/oauth2"
)
func main(){
		conf := &oauth2.Config{
		ClientID:     "cli",            // also known as client key sometimes
		//ClientSecret: "___________________________", // also known as secret key
		Scopes:       []string{"offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://auth.infra.whiteblock.io/auth/realms/wb/protocol/openid-connect/auth",
			TokenURL: "https://auth.infra.whiteblock.io/auth/realms/wb",
		},
	}

	client, err := oauth2ns.AuthenticateUser(conf)
	if err != nil {
		log.Fatal(err)
	}

	// use client.Get / client.Post for further requests, the token will automatically be there
	_, _ = client.Get("/auth-protected-path")
}