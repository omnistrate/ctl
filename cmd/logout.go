package cmd

import (
	"github.com/omnistrate/ctl/utils"

	"github.com/omnistrate/ctl/config"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:          "logout",
	Short:        "Logout from the Omnistrate platform",
	Long:         `The logout command is used to log out from the Omnistrate platform.`,
	Example:      `  omnistrate-ctl logout`,
	RunE:         runLogout,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	err := config.RemoveAuthConfig()
	if err != nil {
		utils.PrintError(err)
		return err
	}
	utils.PrintSuccess("Credentials removed")

	return nil
}
