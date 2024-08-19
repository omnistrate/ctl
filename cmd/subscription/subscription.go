package subscription

import (
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "subscription [operation] [flags]",
	Short:        "Manage subscriptions for your services",
	Long:         `This command helps you manage subscriptions for your services.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(listCmd)

	Cmd.Example = utils.CombineSubCmdExamples(Cmd)
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
