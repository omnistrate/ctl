package domain

import (
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "domain [operation] [flags]",
	Short: "Manage Customer Domains for your service.",
	Long: `This command helps you manage the domains for your service.
These domains are used to access your service in the cloud. You can set up custom domains for each environment type, such as 'prod', 'dev', 'canary', 'staging', 'qa'.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(getCmd)

	Cmd.Example = utils.CombineSubCmdExamples(Cmd)
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
