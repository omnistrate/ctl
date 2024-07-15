package account

import (
	"context"
	"fmt"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	accountconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/account_config_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var (
	accountExample = `
		# List all accounts
		kubectl get accounts

		# List the account with the name 'my-account'
		kubectl get account my-account`
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

	// List aws accounts
	listRes, err := listAccounts(token, "all")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	allAccounts := listRes.AccountConfigs

	// Print accounts table if no account name is provided
	if len(args) == 0 {
		utils.PrintSuccess(fmt.Sprintf("Found %d accounts", len(allAccounts)))
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

func listAccounts(token string, cloudProvider string) (*accountconfigapi.ListAccountConfigResult, error) {
	account, err := httpclientwrapper.NewAccountConfig(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := accountconfigapi.ListAccountConfigRequest{
		Token:             token,
		CloudProviderName: accountconfigapi.CloudProvider(cloudProvider),
	}

	res, err := account.ListAccountConfig(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func printTable(accounts []*accountconfigapi.DescribeAccountConfigResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "ID\tName\tCloud Provider\tAccount ID\tStatus")

	for _, account := range accounts {
		var accountID, cloudProvider string
		if account.AwsAccountID != nil {
			accountID = *account.AwsAccountID
			cloudProvider = "AWS"
		} else {
			accountID = fmt.Sprintf("%s(ProjectID: %s)", *account.GcpProjectID, *account.GcpProjectNumber)
			cloudProvider = "GCP"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			account.ID,
			account.Name,
			cloudProvider,
			accountID,
			account.Status)
	}

	w.Flush()
}
