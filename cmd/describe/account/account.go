package account

import (
	"encoding/json"
	"fmt"
	accountconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/account_config_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	accountExample = `
		# Describe the account with the name
		omnistrate-ctl describe account <name>`
)

// AccountCmd represents the describe command
var AccountCmd = &cobra.Command{
	Use:          "account <name>",
	Short:        "Describe account",
	Long:         `The describe account command displays detailed information about the account.`,
	Example:      accountExample,
	RunE:         Run,
	SilenceUsage: true,
}

func init() {
	AccountCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
}

func Run(cmd *cobra.Command, args []string) error {
	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// List all accounts
	listRes, err := dataaccess.ListAccounts(token, "all")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Filter accounts by name
	var account *accountconfigapi.DescribeAccountConfigResult
	var found bool
	for _, a := range listRes.AccountConfigs {
		if a.Name == args[0] {
			account = a
			found = true
			break
		}
	}

	if !found {
		utils.PrintError(errors.New("account not found"))
		return nil
	}

	// Print account details
	data, err := json.MarshalIndent(account, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	fmt.Println(string(data))

	return nil
}
