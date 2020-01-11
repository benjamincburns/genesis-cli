package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/whiteblock/genesis-cli/pkg/auth"
)

type SearchParams struct {
	Date  string `json:"date"`
	OrgID string `json:"orgID"`
}

type Item struct {
	Text      string `json:"text"`
	Container string `json:"container"`
	Image     string `json:"image"`
	Tag       string `json:"tag"`
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

	orgID, err := getOrgID(orgNameOrId)
	if err != nil {
		return nil, err
	}

	searchParams := SearchParams{
		OrgID: orgID,
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
