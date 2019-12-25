package org

import (
	"encoding/json"

	"github.com/whiteblock/genesis-cli/pkg/config"

	log "github.com/sirupsen/logrus"
)

var conf = config.NewConfig()

// Get returns the stored org name / id if it is given an empty string
// It will also store org if it is not an empty string for use next time
func Get(org string) string {
	if org != "" {
		set(org)
		return org
	}
	return get()
}

func get() string {
	if !conf.UserDir.Exists(conf.OrgFile) {
		return ""
	}

	data, err := conf.UserDir.ReadFile(conf.OrgFile)
	if err != nil {
		log.WithField("error", err).Debug("couldn't read org file")
		return ""
	}

	var org string

	err = json.Unmarshal(data, &org)
	if err != nil {
		log.WithField("error", err).Debug("couldn't parse org name")
		return ""
	}

	return org
}

func set(org string) {
	data, err := json.Marshal(org)
	if err != nil {
		log.WithField("error", err).Debug("couldn't marshal the org name")
	}
	err = conf.UserDir.WriteFile(conf.OrgFile, data)
	if err != nil {
		log.WithField("error", err).Debug("couldn't write to the org file")
	}
}
