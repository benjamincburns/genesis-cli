package util

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func pprint(subj string, attr ...color.Attribute) {
	if viper.GetBool("no-colors") {
		fmt.Println(subj)
		return
	}
	fmt.Println(color.New(attr...).Sprint(subj))
}

func Printf(format string, a ...interface{}) {
	Print(fmt.Sprintf(format, a...))
}

func Print(i interface{}) {
	pprint(fmt.Sprint(i), color.FgHiWhite)
}

func AuthPrintf(format string, a ...interface{}) {
	AuthPrint(fmt.Sprintf(format, a...))
}

func AuthPrint(i interface{}) {
	pprint(fmt.Sprint(i), color.FgCyan)
}
