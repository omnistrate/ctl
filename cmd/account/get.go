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
	getExample = `  # Get all accounts
  omnistrate-ctl account get

  # Get account with name
  omnistrate-ctl account get <name>

  # Get multiple accounts
  omnistrate-ctl account get <name1> <name2> <name3>

  # Get account with ID
  omnistrate-ctl account get <id> --id

  # Get multiple accounts with IDs
  omnistrate-ctl account get <id1> <id2> <id3> --id`
)

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Display one or more accounts",
	Long:    `The account get command displays basic information about one or more accounts.`,
	Example: getExample,
	RunE:    runGet,
	PostRun: func(cmd *cobra.Command, args []string) {
		dataaccess.AskVerifyAccountIfAny()
	},
	SilenceUsage: true,
}

func init() {
	getCmd.Flags().Bool("id", false, "Specify account ID instead of name")
}

func runGet(cmd *cobra.Command, args []string) error {
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
	if ID {
		for _, id := range args {
			var account *accountconfigapi.DescribeAccountConfigResult
			account, err = dataaccess.DescribeAccount(token, id)
			if err != nil {
				utils.PrintError(err)
				return err
			}
			accounts = append(accounts, account)
		}
	} else {
		// List all accounts
		var listRes *accountconfigapi.ListAccountConfigResult
		listRes, err = dataaccess.ListAccounts(token, "all")
		if err != nil {
			utils.PrintError(err)
			return err
		}
		allAccounts := listRes.AccountConfigs

		// Print accounts table if no account name is provided
		if len(args) == 0 {
			utils.PrintSuccess(fmt.Sprintf("%d account(s) found", len(allAccounts)))
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
		for _, name := range args {
			account, ok := accountMap[name]
			if !ok {
				utils.PrintError(fmt.Errorf("account '%s' not found", name))
				continue
			}
			accounts = append(accounts, account)
		}
	}

	printTable(accounts)

	return nil
}

func printTable(accounts []*accountconfigapi.DescribeAccountConfigResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	_, err := fmt.Fprintln(w, "ID\tName\tStatus\tCloud Provider\tTarget Account ID")
	if err != nil {
		return
	}

	for _, account := range accounts {
		var targetAccountID, cloudProvider string
		if account.AwsAccountID != nil {
			targetAccountID = *account.AwsAccountID
			cloudProvider = "AWS"
		} else {
			targetAccountID = fmt.Sprintf("%s(ProjectID: %s)", *account.GcpProjectID, *account.GcpProjectNumber)
			cloudProvider = "GCP"
		}

		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			account.ID,
			account.Name,
			account.Status,
			cloudProvider,
			targetAccountID)
		if err != nil {
			return
		}
	}

	err = w.Flush()
	if err != nil {
		utils.PrintError(err)
	}
}
