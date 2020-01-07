package util

import (
	"fmt"
	"os"

	"github.com/whiteblock/genesis-cli/pkg/message"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	NoMaxArgs = -1
)

/**
 * Unify error messages through function calls
 */

func CheckArguments(cmd *cobra.Command, args []string, min int, max int) {
	if min == max && len(args) != min {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if min == 1 {
			plural = ""
		}
		ErrorFatalf("Invalid number of arguments. "+
			"Expected exactly %d argument%s. Given %d.", min, plural, len(args))
	}
	if len(args) < min {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if min == 1 {
			plural = ""
		}
		ErrorFatalf("Missing arguments. "+
			"Expected atleast %d argument%s. Given %d.", min, plural, len(args))
	}
	if max != NoMaxArgs && len(args) > max {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if max == 1 {
			plural = ""
		}
		ErrorFatalf("Too many arguments. "+
			"Expected atmost %d argument%s. Given %d.", max, plural, len(args))
	}
}

func InvalidArgument(arg string) {
	Errorf("Invalid argument given: %s.", arg)
}

func InvalidInteger(name string, value string, fatal bool) {
	Errorf("Invalid integer, given \"%s\" for %s.", value, name)
	if fatal {
		os.Exit(1)
	}
}

func CheckIntegerBounds(cmd *cobra.Command, name string, val int, min int, max int) {
	if val < min {
		ErrorFatalf("The value given for %s, %d cannot be less than %d.", name, val, min)
	} else if val > max {
		ErrorFatalf("The value given for %s, %d cannot be greater than %d.", name, val, max)
	}
}

func MalformedUsageError(cmd *cobra.Command, err interface{}) {
	pprint(cmd.UsageString())
	ErrorFatal(err)
}

func FlagNotProvidedError(cmd *cobra.Command, flagName string) {
	pprint(cmd.UsageString())
	ErrorFatalf(`missing required flag: "%s"`, flagName)
}

func ErrorFatalf(base string, args ...interface{}) {
	ErrorFatal(fmt.Sprintf(base, args...))
}

func ErrorFatal(err interface{}) {
	Error(err)
	Print(message.FatalErrorMessage)
	os.Exit(1)
}

func Error(err interface{}) {
	pprintln(fmt.Sprint(err), color.FgRed)
}

func Errorf(base string, args ...interface{}) {
	Error(fmt.Sprintf(base, args...))
}
