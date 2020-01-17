package internal

import (
	"os"

	"github.com/whiteblock/genesis-cli/pkg/message"
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"
	"github.com/whiteblock/genesis-cli/pkg/validate"
)

func Lint(data []byte) {
	res, err := validate.Schema(data)
	if err != nil {
		util.ErrorFatal(err)
	}
	if !res.Valid() {
		for _, schemaErr := range res.Errors() {
			util.Error(schemaErr.String())
		}
		os.Exit(1)
	}
	_, _, err = service.ProcessDefinitionFromBytes(data)
	if err != nil {
		util.ErrorFatal(err)
	}

	util.Print(message.FilePassedValidation)
}
