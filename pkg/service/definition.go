package service

import (
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
	"github.com/whiteblock/definition"
	"github.com/whiteblock/definition/command"
)

func ProcessDefinitionFromFile(filename string) (tests []command.Test, def definition.Definition, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	return ProcessDefinition(data)
}

func ProcessDefinition(data []byte) (tests []command.Test, def definition.Definition, err error) {

	v := viper.New()
	v.Set("verbosity", "FATAL")
	definition.ConfigureGlobalFromViper(v) //shush library

	def, err = definition.SchemaYAML(data)
	if err != nil {
		return
	}

	tests, err = definition.GetTests(def)
	return
}
