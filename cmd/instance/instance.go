package instance

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:          "instance [operation] [flags]",
	Short:        "Manage Instance deployment for your service using this command.",
	Long:         `This command helps you manage the deployment of your service instances.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(describeCmd)

	Cmd.Example = instanceExample()

	Cmd.Args = cobra.MinimumNArgs(1)
}

func instanceExample() (example string) {
	for _, cmd := range Cmd.Commands() {
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
