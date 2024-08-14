package account

import (
	"encoding/json"
	"fmt"
	accountconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/account_config_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"slices"
)

const (
	accountExample = `  # Describe account with name
  omnistrate-ctl describe account <name>

  # Describe account with ID
  omnistrate-ctl describe account <id> --id
  
  # Describe multiple accounts with names
  omnistrate-ctl describe account <name1> <name2> <name3>

  # Describe multiple accounts with IDs
  omnistrate-ctl describe account <id1> <id2> <id3> --id`
)

var AccountCmd = &cobra.Command{
	Use:     "account <name>",
	Short:   "Display details for one or more accounts",
	Long:    "Display detailed information about the account by specifying the account name or ID.",
	Example: accountExample,
	RunE:    Run,
	PostRun: func(cmd *cobra.Command, args []string) {
		dataaccess.AskVerifyAccountIfAny()
	},
	SilenceUsage: true,
}

func init() {
	AccountCmd.Args = cobra.MinimumNArgs(1) // Require at least one argument

	AccountCmd.Flags().Bool("id", false, "Specify account ID instead of name")
}

func Run(cmd *cobra.Command, args []string) error {
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

	var accounts []*accountconfigapi.DescribeAccountConfigResult
	for _, name := range args {
		if ID {
			account, err := dataaccess.DescribeAccount(name, token)
			if err != nil {
				utils.PrintError(err)
				return err
			}
			accounts = append(accounts, account)
		} else {
			// List all accounts
			listRes, err := dataaccess.ListAccounts(token, "all")
			if err != nil {
				utils.PrintError(err)
				return err
			}

			// Filter accounts by name
			var found bool
			for _, a := range listRes.AccountConfigs {
				if slices.Contains(args, a.Name) {
					accounts = append(accounts, a)
					found = true
					break
				}
			}

			if !found {
				utils.PrintError(errors.New("account not found: " + name))
				return nil
			}
		}
	}

	// Print account details
	for _, account := range accounts {
		data, err := json.MarshalIndent(account, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}
		fmt.Println(string(data))
	}

	return nil
}
