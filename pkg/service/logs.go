package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/message"
	organization "github.com/whiteblock/genesis-cli/pkg/org"
	"io/ioutil"
	"time"
)

type SearchParams struct {
	Date  string `json:"date"`
	OrgID string `json:"orgID"`
}

type Item struct {
	Text      string `json:"text"`
	Container string `json:"container"`
	Image     string `json:"container"`
	Tag       string `json:"container"`
}

type LogItem struct {
	Timestamp string `json:"timestamp"`
	Message   Item   `json:"item"`
}

func GetLogs(orgNameOrId string) ([]LogItem, error) {
	client, err := auth.GetClient()
	if err != nil {
		return nil, err
	}

	if orgNameOrId == "" {
		orgNameOrId = conf.OrgID
	}
	org, err := organization.Get(orgNameOrId, client)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf(message.MissingOrgID)
	}
	if org.ID == "" {
		return nil, fmt.Errorf(message.MissingOrgID)
	}

	searchParams := SearchParams{
		OrgID: org.ID,
		Date:  time.Now().String(),
	}

	postData, err := json.Marshal(searchParams)
	if err != nil {
		return nil, err
	}

	dest := conf.APIEndpoint() + conf.LogsURI
	resp, err := client.Post(dest, "application/json", bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(data))
	}

	var items []LogItem

	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}

	return items, err
}
