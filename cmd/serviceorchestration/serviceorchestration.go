package serviceorchestration

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "service-orchestration [operation] [flags]",
	Short:        "Manage Service Orchestration Deployments across services",
	Long:         `This command helps you manage orchestration of deployment across multiple services.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(describeCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(listCmd)
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
