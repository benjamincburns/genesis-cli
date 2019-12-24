package util

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func Printf(format string, a ...interface{}) {
	Print(fmt.Sprintf(format, a...))
}

func Print(i interface{}) {
	out := fmt.Sprint(i)

	useColor := !viper.GetBool("no-colors")
	if useColor {
		out = color.WhiteString(out)
	}

	fmt.Println(out)
}
