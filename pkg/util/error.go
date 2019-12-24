package util

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		PrintErrorFatalf("Invalid number of arguments. "+
			"Expected exactly %d argument%s. Given %d.", min, plural, len(args))
	}
	if len(args) < min {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if min == 1 {
			plural = ""
		}
		PrintErrorFatalf("Missing arguments. "+
			"Expected atleast %d argument%s. Given %d.", min, plural, len(args))
	}
	if max != NoMaxArgs && len(args) > max {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if max == 1 {
			plural = ""
		}
		PrintErrorFatalf("Too many arguments. "+
			"Expected atmost %d argument%s. Given %d.", max, plural, len(args))
	}
}

func InvalidArgument(arg string) {
	PrintErrorf("Invalid argument given: %s.", arg)
}

func InvalidInteger(name string, value string, fatal bool) {
	PrintErrorf("Invalid integer, given \"%s\" for %s.", value, name)
	if fatal {
		os.Exit(1)
	}
}

func CheckIntegerBounds(cmd *cobra.Command, name string, val int, min int, max int) {
	if val < min {
		PrintErrorFatalf("The value given for %s, %d cannot be less than %d.", name, val, min)
	} else if val > max {
		PrintErrorFatalf("The value given for %s, %d cannot be greater than %d.", name, val, max)
	}
}

func MalformedUsageError(cmd *cobra.Command, err interface{}) {
	fmt.Println(cmd.UsageString())
	PrintErrorFatal(err)
}

func FlagNotProvidedError(cmd *cobra.Command, flagName string) {
	fmt.Println(cmd.UsageString())
	PrintErrorFatalf(`missing required flag: "%s"`, flagName)
}

func PrintErrorFatalf(base string, args ...interface{}) {
	PrintErrorFatal(fmt.Sprintf(base, args...))
}

func PrintErrorFatal(err interface{}) {
	PrintError(err)
	Print("If you believe this is a bug, please file a bug report")
	os.Exit(1)
}

func PrintError(err interface{}) {
	out := fmt.Sprint(err)

	useColor := !viper.GetBool("no-colors")
	if useColor {
		out = color.RedString(out)
	}
	fmt.Println(out)
}

func PrintErrorf(base string, args ...interface{}) {
	PrintError(fmt.Sprintf(base, args...))
}
