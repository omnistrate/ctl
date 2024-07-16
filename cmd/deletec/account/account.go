package account

import (
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"slices"
	"strings"
)

var (
	accountExample = `  # Delete account with name
  omnistrate-ctl delete account <name>

  # Delete account with ID
  omnistrate-ctl delete account <id> --id

  # Delete multiple accounts with names
  omnistrate-ctl delete account <name1> <name2> <name3>

  # Delete multiple accounts with IDs
  omnistrate-ctl delete account <id1> <id2> <id3> --id`
)

// AccountCmd represents the delete command
var AccountCmd = &cobra.Command{
	Use:          "account <name>",
	Short:        "Delete one or more accounts",
	Long:         `Delete account with name or ID. Use --id to specify ID. If not specified, name is assumed. If multiple accounts are found with the same name, all of them will be deleted.`,
	Example:      accountExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	AccountCmd.Args = cobra.MinimumNArgs(1) // Require at least one argument

	AccountCmd.Flags().Bool("id", false, "Specify account ID instead of name")
}

func run(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var ID bool
	ID, err = cmd.Flags().GetBool("id")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	accountIDs := make([]string, 0)
	if ID {
		accountIDs = args
	} else {
		// List accounts
		listRes, err := dataaccess.ListAccounts(token, "all")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		// Filter accounts by name
		found := make(map[string]int)
		for _, name := range args {
			found[name] = 0
		}

		for _, s := range listRes.AccountConfigs {
			if slices.Contains(args, s.Name) {
				accountIDs = append(accountIDs, string(s.ID))
				found[s.Name] += 1
			}
		}

		accountsNotFound := make([]string, 0)
		for name, count := range found {
			if count == 0 {
				accountsNotFound = append(accountsNotFound, name)
			}
		}

		if len(accountsNotFound) > 0 {
			err = errors.New("account(s) not found: " + strings.Join(accountsNotFound, ", "))
			utils.PrintError(err)
			return err
		}

		for name, count := range found {
			if count > 1 {
				utils.PrintWarning("Multiple accounts found with name: " + name + ". Deleting all of them.")
			}
		}
	}

	// Delete account
	for _, accountID := range accountIDs {
		err = dataaccess.DeleteAccount(accountID, token)
		if err != nil {
			utils.PrintError(err)
			return err
		}
	}

	utils.PrintSuccess("Account(s) deleted successfully")

	return nil
}
