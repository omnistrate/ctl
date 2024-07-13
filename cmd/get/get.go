package get

import (
	"github.com/omnistrate/ctl/cmd/get/account"
	"github.com/omnistrate/ctl/cmd/get/service"
	"github.com/spf13/cobra"
)

var (
	getLong = `
		Display one or many objects.

		Prints a table of the most important information about the specified objects.`
)

var GetCmd = &cobra.Command{
	Use:          "get [object] [name] [flags]",
	Short:        "Display one or many objects",
	Long:         getLong,
	Run:          runGet,
	SilenceUsage: true,
}

func init() {
	GetCmd.AddCommand(service.ServiceCmd)
	GetCmd.AddCommand(account.AccountCmd)

	GetCmd.Example = getExample()

	GetCmd.Args = cobra.MinimumNArgs(1)
}

func getExample() (example string) {
	for _, cmd := range GetCmd.Commands() {
		example += cmd.Example + "\n"
	}
	return
}

func runGet(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
