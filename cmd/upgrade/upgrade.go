package upgrade

import (
	"github.com/spf13/cobra"
)

const (
	upgradeLong = ``

	upgradeExample = `  # Upgrade instances to a specific version
  omnistrate-ctl upgrade <instance1> <instance2> --version 1.2.3

  # Upgrade instances to the latest version
  omnistrate-ctl upgrade <instance1> <instance2> --version latest`
)

var version string

var UpgradeCmd = &cobra.Command{
	Use:          "upgrade <instance> [--version VERSION]",
	Short:        "Upgrade instance to a newer version or an older version.",
	Long:         upgradeLong,
	Example:      upgradeExample,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	UpgradeCmd.Args = cobra.MinimumNArgs(1)

	UpgradeCmd.Flags().StringVarP(&version, "version", "", "", "Specify the version number to upgrade to. Use 'latest' to upgrade to the latest version.")
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}
