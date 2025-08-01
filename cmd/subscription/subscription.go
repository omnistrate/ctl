package subscription

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "subscription [operation] [flags]",
	Short:        "Manage Customer Subscriptions for your service",
	Long:         `This command helps you manage Customer Subscriptions for your service.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(describeCmd)

}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
