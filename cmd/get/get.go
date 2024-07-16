package get

import (
	getaccount "github.com/omnistrate/ctl/cmd/get/account"
	getservice "github.com/omnistrate/ctl/cmd/get/service"
	"github.com/spf13/cobra"
)

var (
	getLong = ``
)

var GetCmd = &cobra.Command{
	Use:          "get [type] [name] [flags]",
	Short:        "Display one or many objects with a table, only the most important information will be displayed.",
	Long:         getLong,
	Run:          runGet,
	SilenceUsage: true,
}

func init() {
	GetCmd.AddCommand(getservice.ServiceCmd)
	GetCmd.AddCommand(getaccount.AccountCmd)

	GetCmd.Example = getExample()

	GetCmd.Args = cobra.MinimumNArgs(1)
}

func getExample() (example string) {
	for _, cmd := range GetCmd.Commands() {
		example += cmd.Example + "\n\n"
	}
	return
}

func runGet(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
