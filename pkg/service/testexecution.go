package service

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/parser"
)

var conf = config.NewConfig()

func TestExecute(filePath string, org string) (string, error) {
	client, err := auth.GetClient()
	if err != nil {
		return "", err
	}

	if org == "" {
		org = conf.OrgID
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
		return string(data), fmt.Errorf("got back a %s code", res.Status)
	}
	return string(data), nil
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
