package util

import (
	"github.com/mattn/go-isatty"
	"os"
)

func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsTerminal(os.Stdout.Fd())
}
