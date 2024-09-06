package instance

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "instance [operation] [flags]",
	Short:        "Manage Instance Deployments for your service",
	Long:         `This command helps you manage the deployment of your service instances.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(describeCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(startCmd)
	Cmd.AddCommand(stopCmd)
	Cmd.AddCommand(restartCmd)
	Cmd.AddCommand(updateCmd) // Hidden (deprecated)
	Cmd.AddCommand(modifyCmd)
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
