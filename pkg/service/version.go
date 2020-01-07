package service

import (
	"io/ioutil"
	"net/http"
)

func CheckForUpdates(version string) <-chan bool {
	out := make(chan bool)
	go func(version string) {
		resp, err := http.Get(conf.VersionLocation)
		if err != nil {
			out <- false
			return
		}
		if resp.StatusCode != 200 {
			out <- false
			return
		}
		vers, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			out <- false
			return
		}
		if len(vers) == 0 {
			out <- false
			return
		}
		out <- version != string(vers)
	}(version)
	return out
}
