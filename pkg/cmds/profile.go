package cmds

import (
	"io/ioutil"
	"os"

	"github.com/whiteblock/genesis-cli/pkg/auth"
	organization "github.com/whiteblock/genesis-cli/pkg/org"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
	"github.com/whiteblock/utility/common"
)

func getOrgId(cmd *cobra.Command) string {
	out := util.GetStringFlagValue(cmd, "org")
	if out == "" {
		return organization.GetDefaultOrgID()
	}
	orgInfo, err := organization.GetOrgInfo(out)
	if err != nil {
		util.ErrorFatal(err)
	}
	return orgInfo.ID
}

var profileCmd = &cobra.Command{
	Use:     "profile",
	Short:   "profile",
	Long:    `profile`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		var out common.OrganizationProfile
		err := auth.Get(conf.GetOrgProfileURL(getOrgId(cmd)), &out)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.PrintS(0, out)
	},
}

var profileSetCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var profileCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := auth.Post(conf.CreateOrgProfileURL(getOrgId(cmd)), nil)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print("success")
	},
}

var profileSetBodyCmd = &cobra.Command{
	Use:   "body [text]",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		var data string

		if len(args) > 0 {
			data = args[0]
		} else {
			fileName := util.GetStringFlagValue(cmd, "file")
			f, err := os.Open(fileName)
			if err != nil {
				util.ErrorFatal(err)
			}
			defer f.Close()
			raw, err := ioutil.ReadAll(f)
			if err != nil {
				util.ErrorFatal(err)
			}
			data = string(raw)
		}
		_, err := auth.Put(conf.UpdateOrgProfileURL(getOrgId(cmd)), [][]string{
			{"body", data},
		})
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print("success")
	},
}

var profileSetEmailCTACmd = &cobra.Command{
	Use:   "email-cta <text> <email>",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		_, err := auth.Put(conf.UpdateOrgProfileURL(getOrgId(cmd)), [][]string{
			{"email_text", args[0]},
			{"email_link", "mailto:" + args[1]},
		})
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print("success")
	},
}

var profileSetWebsiteCTACmd = &cobra.Command{
	Use:   "website-cta <text> <link>",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		_, err := auth.Put(conf.UpdateOrgProfileURL(getOrgId(cmd)), [][]string{
			{"website_text", args[0]},
			{"website_link", args[1]},
		})
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print("success")
	},
}

var profileSetMainCTACmd = &cobra.Command{
	Use:   "main-cta <text> <link>",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)

		toUpdate := [][]string{}

		if cmd.Flags().Changed("icon") {
			toUpdate = append(toUpdate, []string{"main_cta_icon", util.GetStringFlagValue(cmd, "icon")})
		}
		if cmd.Flags().Changed("text") {
			toUpdate = append(toUpdate, []string{"main_cta_text", util.GetStringFlagValue(cmd, "text")})
		}
		if cmd.Flags().Changed("link") {
			toUpdate = append(toUpdate, []string{"main_cta_link", util.GetStringFlagValue(cmd, "link")})
		}

		_, err := auth.Put(conf.UpdateOrgProfileURL(getOrgId(cmd)), toUpdate)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print("success")
	},
}

func init() {
	profileCmd.PersistentFlags().StringP("org", "o", "", "organization")
	profileSetBodyCmd.Flags().StringP("file", "f", "", "set the body from a text file")
	profileSetMainCTACmd.Flags().String("icon", "", "icon")
	profileSetMainCTACmd.Flags().String("text", "", "text")
	profileSetMainCTACmd.Flags().String("link", "", "link")
	profileSetCmd.AddCommand(profileSetBodyCmd)
	profileSetCmd.AddCommand(profileSetEmailCTACmd)
	profileSetCmd.AddCommand(profileSetWebsiteCTACmd)
	profileSetCmd.AddCommand(profileSetMainCTACmd)

	profileCmd.AddCommand(profileSetCmd)
	profileCmd.AddCommand(profileCreateCmd)
	rootCmd.AddCommand(profileCmd)
}
