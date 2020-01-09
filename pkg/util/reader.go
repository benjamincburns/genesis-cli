package util

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

func ReadFileOrRemote(file string) ([]byte, error) {
	rdr, err := FileReader(file)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()
	res, err := ioutil.ReadAll(rdr)
	return res, err
}

func FileReader(file string) (io.ReadCloser, error) {
	r, err := os.Open(file)
	if err == nil {
		return r, nil
	}
	_, err2 := url.Parse(file)
	if err2 != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 1 * time.Minute}
	resp, err := client.Get(file)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func ReadInputFile(basedir string, filename string) (io.ReadCloser, error) {
	file := filename
	if basedir != "" {
		file = path.Join(basedir, filename)
	}
	return FileReader(file)
}
