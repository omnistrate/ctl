package account

import (
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	accountExample = `
		# Delete the account with name
		omnistrate-ctl delete account <name>`
)

// AccountCmd represents the delete command
var AccountCmd = &cobra.Command{
	Use:          "account <name>",
	Short:        "Delete a account",
	Long:         ``,
	Example:      accountExample,
	RunE:         run,
	SilenceUsage: true,
}

func run(cmd *cobra.Command, args []string) error {
	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// List accounts
	listRes, err := dataaccess.ListAccounts(token, "all")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Filter accounts by name
	var accountID string
	var found bool
	for _, s := range listRes.AccountConfigs {
		if s.Name == args[0] {
			accountID = string(s.ID)
			found = true
			break
		}
	}

	if !found {
		utils.PrintError(errors.New("account not found"))
		return nil
	}

	// Delete account
	err = dataaccess.DeleteAccount(accountID, token)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	utils.PrintSuccess("Account deleted successfully")

	return nil
}
