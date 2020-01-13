package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	//log "github.com/sirupsen/logrus"
)

func ReadInputFile(basedir string, filename string) (io.ReadCloser, error) {
	file := filename
	if basedir != "" {
		file = filepath.Join(basedir, filename)
	}
	r, err := os.Open(file)
	if err == nil {
		return r, nil
	}
	if !strings.Contains(filename, "http") {
		return nil, fmt.Errorf(`could not read file "%s" :%s`, file, err.Error())
	}

	_, err2 := url.Parse(filename)
	if err2 != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 1 * time.Minute}
	resp, err := client.Get(filename)
	if err == nil {
		return resp.Body, nil
	}
	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf(string(data))
	}
	return nil, fmt.Errorf(`could not retrieve source-file "%s"`, filename)
}
