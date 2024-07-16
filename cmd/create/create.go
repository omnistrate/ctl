package create

import (
	"github.com/omnistrate/ctl/cmd/create/account"
	"github.com/spf13/cobra"
)

var (
	createLong = ``
)

var CreateCmd = &cobra.Command{
	Use:          "create [object] [name] [flags]",
	Short:        "",
	Long:         createLong,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	CreateCmd.AddCommand(account.AccountCmd)

	CreateCmd.Example = createExample()

	CreateCmd.Args = cobra.MinimumNArgs(1)
}

func createExample() (example string) {
	for _, cmd := range CreateCmd.Commands() {
		example += cmd.Example + "\n"
	}
	return
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
