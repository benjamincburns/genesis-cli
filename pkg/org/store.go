package org

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whiteblock/genesis-cli/pkg/config"
	oauth2ns "github.com/whiteblock/genesis-cli/pkg/oauth2-noserver"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var conf = config.NewConfig()

type Organization struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

// Get returns the stored org information if it is given an empty string
// It will also store information if it is not an empty string for use next time
func Get(orgIdOrName string, client *oauth2ns.AuthorizedClient) (Organization, error) {
	if orgIdOrName != "" {
		return set(orgIdOrName, client)
	}
	org, err := get()
	if err != nil {
		// backwards compatible to just storing org as a string
		return set(legacyGet(), client)
	}
	return org, nil
}

func get() (Organization, error) {
	org := Organization{}
	if !conf.UserDir.Exists(conf.OrgFile) {
		return org, errors.New("no org file")
	}

	data, err := conf.UserDir.ReadFile(conf.OrgFile)
	if err != nil {
		log.WithField("error", err).Debug("couldn't read org file")
		return org, errors.New("couldn't read org file")
	}

	err = json.Unmarshal(data, &org)
	if err != nil {
		log.WithField("error", err).Debug("couldn't parse org")
		return org, err
	}

	return org, nil
}

func legacyGet() string {
	if !conf.UserDir.Exists(conf.OrgFile) {
		return ""
	}

	data, err := conf.UserDir.ReadFile(conf.OrgFile)
	if err != nil {
		log.WithField("error", err).Debug("couldn't read org file")
		return ""
	}

	var org string

	err = json.Unmarshal(data, &org)
	if err != nil {
		log.WithField("error", err).Debug("couldn't parse org name")
		return ""
	}

	return org
}

func set(orgIdOrName string, client *oauth2ns.AuthorizedClient) (Organization, error) {
	org := Organization{}
	dest := conf.APIEndpoint() + fmt.Sprintf(conf.GetOrgURI, orgIdOrName)
	req, err := http.NewRequest("GET", dest, &bytes.Buffer{})
	if err != nil {
		return org, err
	}

	res, err := client.Do(req)
	if res.StatusCode != 200 {
		return org, errors.New("error connecting to backend")
	}

	if err != nil {
		return org, err
	}

	data, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(data, &org)
	if err != nil {
		log.WithField("error", err).Debug("couldn't unmarshal the org")
	}
	err = conf.UserDir.WriteFile(conf.OrgFile, data)
	if err != nil {
		log.WithField("error", err).Debug("couldn't write to the org file")
	}
	return org, nil
}
