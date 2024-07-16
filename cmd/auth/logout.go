package auth

import (
	"github.com/omnistrate/ctl/utils"

	"github.com/omnistrate/ctl/config"
	"github.com/spf13/cobra"
)

// LogoutCmd represents the logout command
var LogoutCmd = &cobra.Command{
	Use:          "logout",
	Short:        "Logout.",
	Long:         `The logout command is used to log out from the Omnistrate platform.`,
	Example:      `  omnistrate-ctl logout`,
	RunE:         runLogout,
	SilenceUsage: true,
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
