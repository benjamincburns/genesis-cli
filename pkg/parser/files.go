package parser

import (
	"fmt"
	"strings"

	"github.com/whiteblock/definition"
	"github.com/whiteblock/definition/pkg/search"
	"github.com/whiteblock/definition/schema"
)

func process(inputFiles []schema.InputFile) []string {
	out := make([]string, len(inputFiles))
	for i, file := range inputFiles {
		out[i] = file.SourcePath
	}
	return out
}

func ExtractFiles(specData []byte) ([]string, error) {

	var root schema.RootSchema
	def, err := definition.SchemaANY(specData)
	if err != nil {
		return nil, err
	}

	root = def.Spec

	max := search.FindServiceMaxCounts(def.Spec)
	files := map[string]bool{}
	for _, service := range root.Services {
		extracted := process(service.InputFiles)
		for _, fileName := range extracted {
			if !strings.Contains(fileName, "$_n") {
				files[fileName] = true
				continue
			}
			for i := int64(0); i < max[service.Name]; i++ {
				files[strings.Replace(fileName, "$_n", fmt.Sprint(i), -1)] = true
			}
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
