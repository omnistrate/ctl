package describe

import (
	"github.com/omnistrate/ctl/cmd/describe/account"
	"github.com/omnistrate/ctl/cmd/describe/service"
	"github.com/spf13/cobra"
)

var (
	describeLong = `
		Describe detailed information about an object.`
)

var DescribeCmd = &cobra.Command{
	Use:          "describe [object] [name] [flags]",
	Short:        "Describe detailed information about an object and output results as JSON to stdout.",
	Long:         describeLong,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	DescribeCmd.AddCommand(service.ServiceCmd)
	DescribeCmd.AddCommand(account.AccountCmd)

	DescribeCmd.Example = describeExample()

	DescribeCmd.Args = cobra.MinimumNArgs(1)
}

func describeExample() (example string) {
	for _, cmd := range DescribeCmd.Commands() {
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
