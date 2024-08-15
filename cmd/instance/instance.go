package instance

import (
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "instance [operation] [flags]",
	Short:        "Manage Instance deployment for your service using this command",
	Long:         `This command helps you manage the deployment of your service instances.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(describeCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(listCmd)

	Cmd.Example = utils.CombineSubCmdExamples(Cmd)
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
