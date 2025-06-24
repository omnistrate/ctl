package environment

import (
	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:          "secret [operation] [flags]",
	Short:        "Manage environment secrets",
	Long:         `This command helps you manage secrets for your service environments.`,
	Run:          runSecret,
	SilenceUsage: true,
}

func init() {
	secretCmd.AddCommand(secretCreateCmd)
	secretCmd.AddCommand(secretListCmd)
	secretCmd.AddCommand(secretDescribeCmd)
	secretCmd.AddCommand(secretUpdateCmd)
	secretCmd.AddCommand(secretDeleteCmd)
}

func runSecret(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}