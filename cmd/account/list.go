package account

import (
	"encoding/json"
	"fmt"
	accountconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/account_config_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	listExample = `# List accounts
omnistrate account list -o=table`
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List cloud provider accounts",
	Long: `This command helps you list cloud provider accounts.
You can filter for specific accounts by using the filter flag.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	listCmd.Flags().StringArrayP("filter", "f", []string{}, "Filter to apply to the list of accounts. E.g.: key1:value1,key2:value2, which filters accounts where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: "+strings.Join(utils.GetSupportedFilterKeys(model.Account{}), ",")+". Check the examples for more details.")
}

func runList(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve command-line flags
	output, _ := cmd.Flags().GetString("output")
	filters, _ := cmd.Flags().GetStringArray("filter")

	// Parse and validate filters
	filterMaps, err := utils.ParseFilters(filters, utils.GetSupportedFilterKeys(model.Account{}))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Ensure user is logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Retrieve accounts and accounts
	listRes, err := dataaccess.ListAccounts(token, "all")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var accountsData []string

	// Process and filter accounts
	for _, account := range listRes.AccountConfigs {
		formattedAccount, err := formatAccount(account)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		match, err := utils.MatchesFilters(formattedAccount, filterMaps)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		data, err := json.MarshalIndent(formattedAccount, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if match {
			accountsData = append(accountsData, string(data))
		}
	}

	// Handle case when no accounts match
	if len(accountsData) == 0 {
		utils.PrintInfo("No accounts found.")
		return nil
	}

	// Format output as requested
	err = utils.PrintTextTableJsonArrayOutput(output, accountsData)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

func formatAccount(account *accountconfigapi.DescribeAccountConfigResult) (model.Account, error) {
	var targetAccountID, cloudProvider string
	if account.AwsAccountID != nil {
		targetAccountID = *account.AwsAccountID
		cloudProvider = "AWS"
	} else {
		targetAccountID = fmt.Sprintf("%s(ProjectID: %s)", *account.GcpProjectID, *account.GcpProjectNumber)
		cloudProvider = "GCP"
	}

	return model.Account{
		ID:              string(account.ID),
		Name:            account.Name,
		Status:          string(account.Status),
		CloudProvider:   cloudProvider,
		TargetAccountID: targetAccountID,
	}, nil
}
