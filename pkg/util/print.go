package util

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func pprintln(subj string, attr ...color.Attribute) {
	if viper.GetBool("no-colors") {
		fmt.Print(subj)
		return
	}
	fmt.Println(color.New(attr...).Sprint(subj))
}

func pprint(subj string, attr ...color.Attribute) {
	if viper.GetBool("no-colors") {
		fmt.Print(subj)
		return
	}
	fmt.Print(color.New(attr...).Sprint(subj))
}

func Printf(format string, a ...interface{}) {
	Print(fmt.Sprintf(format, a...))
}

func Print(i interface{}) {
	pprintln(fmt.Sprint(i), color.FgHiWhite)
}

func PlainPrintf(format string, a ...interface{}) {
	PlainPrint(fmt.Sprintf(format, a...))
}

func PlainPrint(i interface{}) {
	fmt.Print(i)
}

func AuthPrintf(format string, a ...interface{}) {
	AuthPrint(fmt.Sprintf(format, a...))
}

func AuthPrint(i interface{}) {
	pprintln(fmt.Sprint(i), color.FgCyan)
}

func PrintKV(depth int, k interface{}, v interface{}) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "    "
	}
	pprint(indent+fmt.Sprint(k)+": ", color.FgYellow)
	pprintln(fmt.Sprintf("%+v", v))
}
