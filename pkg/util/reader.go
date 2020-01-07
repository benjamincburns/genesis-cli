package util

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
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
	_, err2 := url.Parse(filename)
	if err2 != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 1 * time.Minute}
	resp, err := client.Get(filename)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
