package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/message"
	organization "github.com/whiteblock/genesis-cli/pkg/org"
	"github.com/whiteblock/genesis-cli/pkg/parser"
)

var conf = config.NewConfig()

type Error struct {
	Message string
	Info    []string
}

type Response struct {
	Data interface{} `json:"data"`

	Error *Error   `json:"error,omitempty"`
	Meta  struct{} `json:"meta"`
}

func TestExecute(filePath string, org string) (string, error) {
	client, err := auth.GetClient()
	if err != nil {
		return "", err
	}

	if org == "" {
		org = conf.OrgID
	}
	org = organization.Get(org)

	if org == "" {
		return "", fmt.Errorf(message.MissingOrgID)
	}

	dest := conf.WBHost + fmt.Sprintf(conf.MultipathUploadURI, org)
	if !strings.HasPrefix(dest, "http") {
		dest = "https://" + dest
	}
	req, err := buildRequest(dest, filePath)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return string(data), fmt.Errorf("server responsed with %s", res.Status)
	}

	var resp Response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return string(data), err
	}

	result, ok := resp.Data.(map[string]interface{})
	if !ok {
		return fmt.Sprint(resp), nil
	}

	out := fmt.Sprintf("%v\n", result["message"])
	out += fmt.Sprintf("Definition: %v\n", result["definitionID"])

	return out, nil
}

func buildRequest(dest string, filePath string) (*http.Request, error) {
	b := bytes.Buffer{}
	w := multipart.NewWriter(&b)
	files, err := parser.ExtractFiles(filePath)
	if err != nil {
		return nil, err
	}

	fw, err := w.CreateFormFile("definition", filePath)
	if err != nil {
		return nil, err
	}

	r, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, r)
	if err != nil {
		return nil, err
	}

	for _, fileName := range files {
		fw, err := w.CreateFormFile(fileName, fileName)
		if err != nil {
			return nil, err
		}

		r, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(fw, r)
		if err != nil {
			return nil, err
		}

	}

	w.Close()

	req, err := http.NewRequest("PUT", dest, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}
