package serviceplan

import (
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "service-plan [operation] [flags]",
	Short:        "Manage service plans for your services",
	Long:         `This command helps you manage the service plans for your services.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(releaseCmd)
	Cmd.AddCommand(setDefaultCmd)
	Cmd.AddCommand(describeCmd)
	Cmd.AddCommand(describeVersionCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(listVersionsCmd)

	Cmd.Example = utils.CombineSubCmdExamples(Cmd)
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
