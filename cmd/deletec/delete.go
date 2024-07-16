package deletec

import (
	"github.com/omnistrate/ctl/cmd/deletec/account"
	"github.com/omnistrate/ctl/cmd/deletec/service"
	"github.com/spf13/cobra"
)

var (
	deleteLong = ``
)

var DeleteCmd = &cobra.Command{
	Use:          "delete [object] [name] [flags]",
	Short:        "Delete objects by specifying the object type and name.",
	Long:         deleteLong,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	DeleteCmd.AddCommand(service.ServiceCmd)
	DeleteCmd.AddCommand(account.AccountCmd)

	DeleteCmd.Example = deleteExample()

	DeleteCmd.Args = cobra.MinimumNArgs(1)
}

func deleteExample() (example string) {
	for _, cmd := range DeleteCmd.Commands() {
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
