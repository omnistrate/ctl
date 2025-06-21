package deploymentcell

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "deployment-cell [operation] [flags]",
	Short:        "Manage Deployment Cells",
	Long:         `This command helps you manage Deployment Cells.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(adoptCmd)
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(listCmd)
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}