package util

import (
	"github.com/spf13/cobra"
)

func RequireFlags(cmd *cobra.Command, flags ...string) {
	for _, flag := range flags {
		if !cmd.Flags().Changed(flag) {
			FlagNotProvidedError(cmd, flag)
			return
		}
	}
}

func GetStringFlagValue(cmd *cobra.Command, flag string) string {
	out, err := cmd.Flags().GetString(flag)
	if err != nil {
		out, err = cmd.PersistentFlags().GetString(flag)
		if err != nil {
			ErrorFatal(err)
		}
	}
	return out
}

func GetStringSliceFlagValue(cmd *cobra.Command, flag string) []string {
	out, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		out, err = cmd.PersistentFlags().GetStringSlice(flag)
		if err != nil {
			ErrorFatal(err)
		}
	}
	return out
}

func GetIntFlagValue(cmd *cobra.Command, flag string) int {
	out, err := cmd.Flags().GetInt(flag)
	if err != nil {
		ErrorFatal(err)
	}
	return out
}

func GetInt64FlagValue(cmd *cobra.Command, flag string) string {
	out, err := cmd.Flags().GetInt64(flag)
	if err != nil {
		out, err = cmd.PersistentFlags().GetInt64(flag)
		if err != nil {
			ErrorFatal(err)
		}
	}
	return out
}

func GetFloat64FlagValue(cmd *cobra.Command, flag string) float64 {
	out, err := cmd.Flags().GetFloat64(flag)
	if err != nil {
		ErrorFatal(err)
	}
	return out
}

func GetBoolFlagValue(cmd *cobra.Command, flag string) bool {
	out, err := cmd.Flags().GetBool(flag)
	if err != nil {
		ErrorFatal(err)
	}
	return out
}
