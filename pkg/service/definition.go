package service

import (
	"io/ioutil"
	"os"

	"github.com/whiteblock/definition"
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/schema"

	"github.com/spf13/viper"
)

func init() {
	v := viper.New()
	v.Set("verbosity", "FATAL")
	definition.ConfigureGlobalFromViper(v) //shush library
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

func ReplaceFile(root *schema.RootSchema, find string, replace string) {
	for i := range root.Services {
		for j := range root.Services[i].InputFiles {
			if root.Services[i].InputFiles[j].SourcePath == find {
				root.Services[i].InputFiles[j].SourcePath = replace
			}
		}
	}

	for i := range root.Sidecars {
		for j := range root.Sidecars[i].InputFiles {
			if root.Sidecars[i].InputFiles[j].SourcePath == find {
				root.Sidecars[i].InputFiles[j].SourcePath = replace
			}
		}
	}

	for i := range root.TaskRunners {
		for j := range root.TaskRunners[i].InputFiles {
			if root.TaskRunners[i].InputFiles[j].SourcePath == find {
				root.TaskRunners[i].InputFiles[j].SourcePath = replace
			}
		}
	}
}
