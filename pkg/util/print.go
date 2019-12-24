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
		out = color.HiWhiteString(out)
	}

	fmt.Println(out)
}

func AuthPrintf(format string, a ...interface{}) {
	AuthPrint(fmt.Sprintf(format, a...))
}

func AuthPrint(i interface{}) {
	out := fmt.Sprint(i)

	useColor := !viper.GetBool("no-colors")
	if useColor {
		out = color.CyanString(out)
	}

	fmt.Println(out)
}
