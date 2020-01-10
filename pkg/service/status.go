package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/whiteblock/genesis-cli/pkg/auth"

	"github.com/whiteblock/utility/common"
)

func GetStatus(testID string) (common.Status, error) {
	client, err := auth.GetClient()
	if err != nil {
		return common.Status{}, err
	}

	dest := conf.APIEndpoint() + fmt.Sprintf(conf.StatusURI, testID)
	resp, err := client.Get(dest)
	if err != nil {
		return common.Status{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return common.Status{}, err
	}
	var status common.Status
	return status, json.Unmarshal(data, &status)
}