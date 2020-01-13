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
	"sort"

	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/message"
	organization "github.com/whiteblock/genesis-cli/pkg/org"
	"github.com/whiteblock/genesis-cli/pkg/parser"
	"github.com/whiteblock/genesis-cli/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/whiteblock/definition/schema"
	"github.com/whiteblock/utility/common"
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

func GetMostRecentTest(org string) (common.Test, error) {
	tests, err := GetTests(org)
	if err != nil {
		return common.Test{}, err
	}
	if len(tests) == 0 {
		return common.Test{}, fmt.Errorf("no active tests")
	}
	return tests[len(tests)-1], nil
}

func getOrgID(orgID string) (string, error) {
	client, err := auth.GetClient()
	if err != nil {
		return "", err
	}
	if orgID == "" {
		orgID = conf.OrgID
	}
	org, err := organization.Get(orgID, client)
	if err != nil {
		log.WithField("error", err).Trace("failed to fetch org id ")
		if orgID == "" {
			return orgID, fmt.Errorf(message.MissingOrgID)
		}
		return orgID, nil
	}
	if org.ID == "" && orgID == "" {
		return orgID, fmt.Errorf(message.MissingOrgID)
	}
	return org.ID, nil
}

func GetTests(orgNameOrId string) ([]common.Test, error) {
	client, err := auth.GetClient()
	if err != nil {
		return nil, err
	}

	orgID, err := getOrgID(orgNameOrId)
	if err != nil {
		return nil, err
	}

	dest := conf.APIEndpoint() + fmt.Sprintf(conf.TestsURI, orgID)
	log.WithField("url", dest).Debug("getting tests")
	resp, err := client.Get(dest)
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

	var out []common.Test
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Unix() < out[j].CreatedAt.Unix()
	})
	return out, nil
}

func UploadFiles(filePath string, orgNameOrId string) (string, error) {
	client, err := auth.GetClient()
	if err != nil {
		return "", err
	}

	orgID, err := getOrgID(orgNameOrId)
	if err != nil {
		return "", err
	}

	dest := conf.APIEndpoint() + fmt.Sprintf(conf.MultipathUploadURI, orgID)

	req, err := buildRequest(dest, filePath)
	if err != nil {
		return "", err
	}
	log.WithField("url", dest).Debug("uploading files")
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
	log.WithField("response", string(data)).Trace("got a response from the server")
	var resp Response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return string(data), err
	}

	result, ok := resp.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("%v", resp.Data)
	}
	return fmt.Sprint(result["definitionID"]), nil
}

func RunTest(defFile string, orgNameOrId string, definitionID string, dns []string) (out []string, err error) {
	orgID, err := getOrgID(orgNameOrId)
	if err != nil {
		return
	}

	client, err := auth.GetClient()
	if err != nil {
		return
	}

	r, err := util.ReadInputFile("", defFile)
	if err != nil {
		return nil, err
	}

	dest := conf.APIEndpoint() + fmt.Sprintf(conf.RunTestURI, orgID, definitionID)
	log.WithField("url", dest).Debug("running test")
	req, err := http.NewRequest("POST", dest, r)
	if err != nil {
		return nil, err
	}
	for i := range dns {
		req.Header.Set(DNSHeader, dns[i])
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf(string(data))
	}
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

func StopTest(id string, isDef bool) error {
	client, err := auth.GetClient()
	if err != nil {
		return err
	}
	dest := func() string {
		if isDef {
			return conf.APIEndpoint() + fmt.Sprintf(conf.StopDefURI, id)
		}
		return conf.APIEndpoint() + fmt.Sprintf(conf.StopTestURI, id)
	}()
	resp, err := client.Post(dest, "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf(string(data))
	}
	return nil
}

func buildRequest(dest string, filePath string) (*http.Request, error) {
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
	return req, nil
}
