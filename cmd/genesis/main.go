package main

import (
	"github.com/whiteblock/genesis-cli/pkg/cmds"
)

var (
	buildTime  string
	commitHash string
)

func main() {
	cmds.Execute(buildTime, commitHash)
}
