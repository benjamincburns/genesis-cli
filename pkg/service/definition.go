package service

import (
	"io/ioutil"
	"os"

	"github.com/whiteblock/definition"
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/pkg/entity"

	"github.com/spf13/viper"
)

func init() {
	v := viper.New()
	v.Set("verbosity", conf.Verbosity)
	definition.ConfigureGlobalFromViper(v)
}

func parseDef(data []byte) (definition.Definition, error) {
	def, err := definition.SchemaYAML(data)
	if err != nil {
		return definition.SchemaJSON(data)
	}
	return def, nil
}

func DefinitionFromFile(filename string) (def definition.Definition, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	return parseDef(data)
}

func ProcessDefinitionFromFile(filename string) (tests []command.Test, def definition.Definition, err error) {
	def, err = DefinitionFromFile(filename)
	if err != nil {
		return
	}
	tests, err = definition.GetTests(def, definition.Meta{})
	return
}

func ProcessDefinitionFromBytes(data []byte) (tests []command.Test, def definition.Definition, err error) {
	def, err = parseDef(data)
	if err != nil {
		return
	}
	tests, err = definition.GetTests(def, definition.Meta{})
	return
}

func DefinitionEnv(filename string) ([]map[string]string, definition.Definition, error) {
	def, err := DefinitionFromFile(filename)
	if err != nil {
		return nil, definition.Definition{}, err
	}
	envs, err := definition.GetEnvs(def)
	return envs, def, err
}

func DefinitionDist(filename string) ([]*entity.ResourceDist, definition.Definition, error) {
	def, err := DefinitionFromFile(filename)
	if err != nil {
		return nil, definition.Definition{}, err
	}
	dist, err := definition.GetDist(def)
	return dist, def, err
}
