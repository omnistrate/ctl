package account

import (
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "account [operation] [flags]",
	Short:        "Manage your Cloud Provider Accounts",
	Long:         `This command helps you manage your cloud provider accounts.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(describeCmd)
	Cmd.AddCommand(deleteCmd)

	Cmd.Example = utils.CombineSubCmdExamples(Cmd)
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
