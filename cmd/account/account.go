package account

import (
	"github.com/spf13/cobra"
)

var AccountCmd = &cobra.Command{
	Use:          "account [operation] [flags]",
	Short:        "Manage your cloud provider accounts.",
	Long:         `This command helps you manage your cloud provider accounts.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	AccountCmd.AddCommand(createCmd)
	AccountCmd.AddCommand(getCmd)
	AccountCmd.AddCommand(describeCmd)
	AccountCmd.AddCommand(deleteCmd)

	AccountCmd.Example = accountExample()
}

func accountExample() (example string) {
	for _, cmd := range AccountCmd.Commands() {
		example += cmd.Example + "\n\n"
	}
	return
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
