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
	"path/filepath"

	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/message"
	organization "github.com/whiteblock/genesis-cli/pkg/org"
	"github.com/whiteblock/genesis-cli/pkg/parser"
	"github.com/whiteblock/genesis-cli/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/whiteblock/definition/schema"
)

var conf = config.NewConfig()

const DNSHeader = "Wbdns"

type Error struct {
	Message string
	Info    []string
}

type Response struct {
	Data interface{} `json:"data"`

	Error *Error   `json:"error,omitempty"`
	Meta  struct{} `json:"meta"`
}

func TestExecute(filePath string, org string, dns []string) (string, []string, error) {
	client, err := auth.GetClient()
	if err != nil {
		return "", nil, err
	}

	if org == "" {
		org = conf.OrgID
	}
	org = organization.Get(org)

	if org == "" {
		return "", nil, fmt.Errorf(message.MissingOrgID)
	}

	dest := conf.APIEndpoint() + fmt.Sprintf(conf.MultipathUploadURI, org)

	req, err := buildRequest(dest, filePath, dns)
	if err != nil {
		return "", nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, err
	}
	if res.StatusCode != http.StatusOK {
		return string(data), nil, fmt.Errorf("server responsed with %s", res.Status)
	}
	log.WithField("response", string(data)).Trace("got a response from the server")
	var resp Response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return string(data), nil, err
	}

	result, ok := resp.Data.(map[string]interface{})
	if !ok {
		return fmt.Sprint(resp), nil, nil
	}

	out := fmt.Sprintf("%v\n", result["message"])
	out += fmt.Sprintf("Definition: %v\n", result["definitionID"])
	ids := []string{}
	if tests, ok := result["testIDs"]; ok {
		for _, test := range tests.([]interface{}) {
			out += fmt.Sprintf("\tTest: %v\n", test)
			ids = append(ids, fmt.Sprintf("%v", test))
		}
	}

	return out, ids, nil
}

func buildRequest(dest string, filePath string, dns []string) (*http.Request, error) {
	b := bytes.Buffer{}
	w := multipart.NewWriter(&b)
	files, err := parser.ExtractFiles(filePath)
	if err != nil {
		return nil, err
	}

	var root schema.RootSchema

	{
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		def, err := parseDef(data)
		if err != nil {
			return nil, err
		}
		root = def.Spec
	}

	basedir := filepath.Dir(filePath)
	readyFiles := map[string]bool{}
	for i, fileName := range files {
		if _, ok := readyFiles[fileName]; ok {
			continue
		}
		readyFiles[fileName] = true
		fw, err := w.CreateFormFile("/f/k"+fmt.Sprint(i), "/f/k"+fmt.Sprint(i))
		if err != nil {
			return nil, err
		}
		ReplaceFile(&root, fileName, "/f/k"+fmt.Sprint(i))
		r, err := util.ReadInputFile(basedir, fileName)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		_, err = io.Copy(fw, r)
		if err != nil {
			return nil, err
		}

	}

	fw, err := w.CreateFormFile("definition", filePath)
	if err != nil {
		return nil, err
	}
	log.WithField("root", fmt.Sprintf("%+v", root)).Trace("resulting spec")
	data, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(data)
	_, err = io.Copy(fw, r)
	if err != nil {
		return nil, err
	}

	w.Close()

	req, err := http.NewRequest("PUT", dest, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	for i := range dns {
		req.Header.Set(DNSHeader, dns[i])
	}
	return req, nil
}
