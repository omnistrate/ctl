package account

import (
	"fmt"
	accountconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/account_config_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var (
	accountExample = `
		# List all accounts
		omnistrate-ctl get accounts

		# List the account with the name
		omnistrate-ctl get account <name>`
)

// AccountCmd represents the describe command
var AccountCmd = &cobra.Command{
	Use:          "account <name>",
	Short:        "Display one or more accounts",
	Long:         `The get account command displays basic information about one or more accounts.`,
	Example:      accountExample,
	RunE:         Run,
	SilenceUsage: true,
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
	allAccounts := listRes.AccountConfigs

	// Print accounts table if no account name is provided
	if len(args) == 0 {
		utils.PrintSuccess(fmt.Sprintf("%d accounts found", len(allAccounts)))
		if len(allAccounts) > 0 {
			printTable(allAccounts)
		}
		return nil
	}

	// Format listRes.Accounts into a map
	accountMap := make(map[string]*accountconfigapi.DescribeAccountConfigResult)
	for _, account := range allAccounts {
		accountMap[account.Name] = account
	}

	// Filter accounts by name
	var accounts []*accountconfigapi.DescribeAccountConfigResult
	for _, name := range args {
		account, ok := accountMap[name]
		if !ok {
			utils.PrintError(fmt.Errorf("account '%s' not found", name))
			continue
		}
		accounts = append(accounts, account)
	}

	// Print accounts table if no account name is provided
	printTable(accounts)

	return nil
}

func printTable(accounts []*accountconfigapi.DescribeAccountConfigResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "ID\tName\tStatus\tCloud Provider\tTarget Account ID")

	for _, account := range accounts {
		var targetAccountID, cloudProvider string
		if account.AwsAccountID != nil {
			targetAccountID = *account.AwsAccountID
			cloudProvider = "AWS"
		} else {
			targetAccountID = fmt.Sprintf("%s(ProjectID: %s)", *account.GcpProjectID, *account.GcpProjectNumber)
			cloudProvider = "GCP"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			account.ID,
			account.Name,
			account.Status,
			cloudProvider,
			targetAccountID)
	}

	w.Flush()
}
