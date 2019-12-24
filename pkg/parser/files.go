package parser

import (
	"io/ioutil"
	"os"

	"github.com/whiteblock/definition/schema"
	"gopkg.in/yaml.v2"
)

func process(inputFiles []schema.InputFile) []string {
	out := make([]string, len(inputFiles))
	for i, file := range inputFiles {
		out[i] = file.SourcePath
	}
	return out
}

func ExtractFiles(yamlFile string) ([]string, error) {
	f, err := os.Open(yamlFile)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var root schema.RootSchema
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}
	files := map[string]bool{}
	for _, service := range root.Services {
		extracted := process(service.InputFiles)
		for _, fileName := range extracted {
			files[fileName] = true
		}
	}

	for _, sidecar := range root.Sidecars {
		extracted := process(sidecar.InputFiles)
		for _, fileName := range extracted {
			files[fileName] = true
		}
	}

	for _, taskrunner := range root.TaskRunners {
		extracted := process(taskrunner.InputFiles)
		for _, fileName := range extracted {
			files[fileName] = true
		}
	}

	out := []string{}
	for fileName := range files {
		out = append(out, fileName)
	}
	return out, nil
}
